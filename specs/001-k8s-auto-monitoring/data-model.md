# Data Model: Kubernetes 자동 모니터링 시스템

**Date**: 2024-12-01  
**Feature**: 001-k8s-auto-monitoring

## Core Entities

### 1. MonitoringNode (모니터링 노드)

**Description**: Kubernetes 클러스터의 개별 노드로, 모니터링 대상으로 관리되는 엔티티

**Attributes**:
- `name` (string, required): 노드 이름 (Kubernetes Node 리소스의 이름)
- `uid` (string, required): 노드의 고유 식별자 (Kubernetes UID)
- `status` (enum, required): 노드 상태
  - `pending`: 모니터링 설정 중
  - `monitoring`: 모니터링 중
  - `removed`: 모니터링 중단됨
  - `error`: 모니터링 설정 실패
- `nodeExporterDeployed` (boolean, required): node-exporter DaemonSet 배포 여부
- `serviceMonitorCreated` (boolean, required): ServiceMonitor 리소스 생성 여부
- `lastDetectedAt` (timestamp, required): 마지막 감지 시간
- `lastUpdatedAt` (timestamp, required): 마지막 업데이트 시간

**Relationships**:
- 1:1 with Kubernetes Node 리소스
- 1:1 with node-exporter DaemonSet
- 1:1 with ServiceMonitor 리소스

**Validation Rules**:
- 노드 이름은 Kubernetes 리소스 이름 규칙을 따라야 함 (RFC 1123)
- UID는 변경 불가능하며 노드의 고유 식별자로 사용

**State Transitions**:
```
[Node Added] → pending → [node-exporter 배포] → monitoring
[Node Removed] → monitoring → removed
[Error] → pending/monitoring → error
```

---

### 2. MonitoringConfiguration (모니터링 구성)

**Description**: Prometheus, Grafana, Alertmanager의 설정 정보 및 모니터링 대상 목록

**Attributes**:
- `prometheusEnabled` (boolean, required): Prometheus 배포 여부
- `grafanaEnabled` (boolean, required): Grafana 배포 여부
- `alertmanagerEnabled` (boolean, required): Alertmanager 배포 여부
- `namespace` (string, required): 모니터링 컴포넌트가 배포될 네임스페이스 (기본값: "monitoring")
- `prometheusRetention` (string, optional): Prometheus 데이터 보관 기간 (기본값: "15d")
- `grafanaAdminPassword` (string, optional): Grafana 관리자 비밀번호 (자동 생성 가능)
- `alertRules` (array, optional): 커스텀 알림 규칙 목록

**Relationships**:
- 1:N with MonitoringNode (여러 노드를 모니터링)
- 1:1 with Prometheus Deployment
- 1:1 with Grafana Deployment
- 1:1 with Alertmanager Deployment

**Validation Rules**:
- 네임스페이스는 Kubernetes 리소스 이름 규칙을 따라야 함
- Prometheus 보관 기간은 유효한 기간 형식이어야 함 (예: "15d", "720h")

---

### 3. MetricData (메트릭 데이터)

**Description**: 노드와 클러스터에서 수집되는 성능 및 상태 정보

**Attributes**:
- `nodeName` (string, required): 메트릭이 수집된 노드 이름
- `timestamp` (timestamp, required): 메트릭 수집 시간
- `cpuUsage` (float, optional): CPU 사용률 (0-100)
- `memoryUsage` (float, optional): 메모리 사용률 (0-100)
- `diskUsage` (float, optional): 디스크 사용률 (0-100)
- `networkIn` (float, optional): 네트워크 수신 트래픽 (bytes/sec)
- `networkOut` (float, optional): 네트워크 송신 트래픽 (bytes/sec)

**Relationships**:
- N:1 with MonitoringNode (한 노드에서 여러 메트릭 데이터)
- 메트릭은 Prometheus에 저장되며, 이 엔티티는 메트릭 스키마를 나타냄

**Storage**:
- Prometheus 시계열 데이터베이스에 저장
- Grafana를 통해 쿼리 및 시각화

---

### 4. AlertRule (알림 규칙)

**Description**: Prometheus 알림 규칙 정의

**Attributes**:
- `name` (string, required): 알림 규칙 이름
- `alert` (string, required): 알림 이름
- `expr` (string, required): PromQL 표현식
- `for` (string, optional): 알림이 발송되기 전 대기 시간 (예: "5m")
- `labels` (map, optional): 알림에 추가할 레이블
- `annotations` (map, optional): 알림 설명 및 메시지
- `severity` (enum, optional): 심각도 (info, warning, critical)

**Relationships**:
- N:1 with MonitoringConfiguration (여러 알림 규칙이 하나의 구성에 속함)

**Default Rules**:
1. **NodeDown**: 노드가 다운된 경우
2. **HighCPUUsage**: CPU 사용률이 80% 이상인 경우
3. **HighMemoryUsage**: 메모리 사용률이 80% 이상인 경우
4. **DiskSpaceLow**: 디스크 공간이 20% 미만인 경우

---

## Kubernetes Resources

### Node (Kubernetes Native)
- **Type**: Core Kubernetes Resource
- **Purpose**: 클러스터의 물리적/가상 노드
- **Watch Target**: Controller가 이 리소스를 Watch하여 노드 추가/제거 감지

### DaemonSet (node-exporter)
- **Type**: Kubernetes Workload
- **Purpose**: 각 노드에 node-exporter Pod 배포
- **Selector**: `app: node-exporter`
- **Template**: node-exporter 컨테이너 이미지

### Service (node-exporter)
- **Type**: Kubernetes Service
- **Purpose**: node-exporter 엔드포인트 노출
- **Port**: 9100 (node-exporter 기본 포트)
- **Selector**: `app: node-exporter`

### ServiceMonitor (Prometheus Operator)
- **Type**: Custom Resource (monitoring.coreos.com/v1)
- **Purpose**: Prometheus가 node-exporter를 스크랩하도록 설정
- **Selector**: node-exporter Service 선택
- **Endpoints**: `/metrics` 경로 스크랩 설정

### PrometheusRule (Prometheus Operator)
- **Type**: Custom Resource (monitoring.coreos.com/v1)
- **Purpose**: Prometheus 알림 규칙 정의
- **Groups**: 알림 규칙 그룹 배열

### ConfigMap (Grafana Dashboards)
- **Type**: Kubernetes ConfigMap
- **Purpose**: Grafana 대시보드 JSON 저장
- **Label**: `grafana_dashboard: "1"` (자동 로드)
- **Data**: 대시보드 JSON 파일

---

## Data Flow

```
1. Node Added Event
   ↓
2. Controller detects new node
   ↓
3. Create node-exporter DaemonSet (if not exists)
   ↓
4. Create ServiceMonitor for the node
   ↓
5. Prometheus Operator detects ServiceMonitor
   ↓
6. Prometheus starts scraping node-exporter
   ↓
7. Metrics stored in Prometheus
   ↓
8. Grafana queries Prometheus for visualization
```

---

## State Management

### Controller State
- Controller는 Kubernetes API를 통해 실제 상태를 확인
- Desired state는 코드 로직에 정의
- Reconcile 루프를 통해 desired state와 actual state 동기화

### Persistent State
- Prometheus: 메트릭 데이터를 PVC에 저장
- Grafana: 대시보드 및 설정을 ConfigMap에 저장
- Alertmanager: 알림 규칙을 PrometheusRule CRD에 저장

### Ephemeral State
- Controller의 메모리 상태 (재시작 시 재구성)
- 노드 감지 이벤트 (이벤트 기반, 상태 저장 불필요)

