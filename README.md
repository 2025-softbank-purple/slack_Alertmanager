# Kubernetes 모니터링 + Slack Alertmanager

Prometheus Operator(kube-prometheus-stack)와 Alertmanager를 배포해 노드 메트릭을 수집하고, Slack 웹훅으로 알림을 보내는 예제입니다. node-exporter로 노드 상태를 수집하며, 알림은 Secret로 주입한 Slack 웹훅 URL을 사용합니다.

## 구성
- Helm 값 파일: `charts/prometheus-stack/values.yaml` (Alertmanager Slack 설정 포함)
- 경보 규칙 및 ServiceMonitor: `configs/prometheus/`
- node-exporter DaemonSet/Service: `configs/node-exporter/`
- 설치/제거 스크립트: `scripts/`
- 명세/테스트 문서: `specs/`, `tests/`

## Alertmanager → Slack 포맷
- Slack 웹훅 URL: Secret `slack-webhook`의 키 `api_url`
- Alertmanager가 `/etc/alertmanager/secrets/slack-webhook/api_url`에서 URL을 읽음
- 제목: `[STATUS] alertname (severity)`
- 본문:
  - `*Where*: ns=<namespace>, pod=<pod>, instance=<instance>`
  - `*What*: <summary>`
  - `*Detail*: <description>` (주석이 있을 때만)

## 선행 조건
- Docker Desktop 실행 중(WSL2 엔진 권장)
- kind 설치 및 PATH 등록 (`kind --version` 동작)
- kubectl, Helm v3 설치

## 로컬 실행(Windows PowerShell)
```powershell
# 1) kind 클러스터 생성
kind create cluster --config kind-multi-node.yaml

# 2) 네임스페이스
kubectl create namespace monitoring

# 3) Slack 웹훅 Secret
kubectl -n monitoring create secret generic slack-webhook --from-literal=api_url='https://hooks.slack.com/services/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'

# 4) Prometheus/Alertmanager 배포
helm upgrade --install prometheus charts/prometheus-stack -n monitoring -f charts/prometheus-stack/values.yaml

# 5) Alertmanager UI 포트포워드(창 유지)
kubectl -n monitoring port-forward svc/alertmanager 9093:9093
```
Alertmanager UI: http://localhost:9093

### 테스트 알람 보내기
```powershell
curl -XPOST -H "Content-Type: application/json" http://localhost:9093/api/v1/alerts -d "[{`"labels`":{`"alertname`":`"TestAlert`",`"severity`":`"warning`",`"instance`":`"test.local`",`"namespace`":`"default`",`"pod`":`"demo-123`"},`"annotations`":{`"summary`":`"Test summary`",`"description`":`"This is a test alert to Slack`"}}]"
```

### 정리
```powershell
helm uninstall prometheus -n monitoring
kubectl delete namespace monitoring
kind delete cluster
```

## 파일 맵
```
charts/prometheus-stack/values.yaml   # Prometheus/Alertmanager 설정 + Slack 수신자
configs/prometheus/alert-rules.yaml   # 노드 알림 규칙(30s 평가, 5m 지속)
configs/prometheus/servicemonitor.yaml# node-exporter 스크레이프 30s
configs/node-exporter/daemonset.yaml  # node-exporter DaemonSet
configs/node-exporter/service.yaml    # node-exporter Service
scripts/install.sh, scripts/uninstall.sh
specs/, tests/                       # 사양/테스트 문서
```

## 트러블슈팅
- Slack 포맷이 안 바뀔 때: 최신 `values.yaml`로 Helm 재배포 후 Alertmanager StatefulSet 재시작.
- 웹훅 교체: Secret `slack-webhook`의 `api_url`만 갱신하면 재배포 없이 반영.
