# Kubernetes ����͸� + Slack Alertmanager

Prometheus Operator(kube-prometheus-stack)�� Alertmanager�� ������ ��� ��Ʈ���� �����ϰ�, Slack �������� �˸��� ������ �����Դϴ�. node-exporter�� ��� ���¸� �����ϸ�, �˸��� Secret�� ������ Slack ���� URL�� ����մϴ�.

## ����
- Helm �� ����: `charts/prometheus-stack/values.yaml` (Alertmanager Slack ���� ����)
- �溸 ��Ģ �� ServiceMonitor: `configs/prometheus/`
- node-exporter DaemonSet/Service: `configs/node-exporter/`
- ��ġ/���� ��ũ��Ʈ: `scripts/`
- ����/�׽�Ʈ ����: `specs/`, `tests/`

## Alertmanager �� Slack ����
- Slack ���� URL: Secret `slack-webhook`�� Ű `api_url`
- Alertmanager�� `/etc/alertmanager/secrets/slack-webhook/api_url`���� URL�� ����
- ����: `[STATUS] alertname (severity)`
- ����:
  - `*Where*: ns=<namespace>, pod=<pod>, instance=<instance>`
  - `*What*: <summary>`
  - `*Detail*: <description>` (�ּ��� ���� ����)

## ���� ����
- Docker Desktop ���� ��(WSL2 ���� ����)
- kind ��ġ �� PATH ��� (`kind --version` ����)
- kubectl, Helm v3 ��ġ

## ���� ����(Windows PowerShell)
```powershell
# 1) kind Ŭ������ ����
kind create cluster --config kind-multi-node.yaml

# 2) ���ӽ����̽�
kubectl create namespace monitoring

# 3) Slack ���� Secret
kubectl -n monitoring create secret generic slack-webhook --from-literal=api_url='YOUR_SLACK_WEBHOOK_URL_HERE'

# 4) Prometheus/Alertmanager ����
helm upgrade --install prometheus charts/prometheus-stack -n monitoring -f charts/prometheus-stack/values.yaml

# 5) Alertmanager UI ��Ʈ������(â ����)
kubectl -n monitoring port-forward svc/alertmanager 9093:9093
```
Alertmanager UI: http://localhost:9093

### �׽�Ʈ �˶� ������
```powershell
curl -XPOST -H "Content-Type: application/json" http://localhost:9093/api/v1/alerts -d "[{`"labels`":{`"alertname`":`"TestAlert`",`"severity`":`"warning`",`"instance`":`"test.local`",`"namespace`":`"default`",`"pod`":`"demo-123`"},`"annotations`":{`"summary`":`"Test summary`",`"description`":`"This is a test alert to Slack`"}}]"
```

### ����
```powershell
helm uninstall prometheus -n monitoring
kubectl delete namespace monitoring
kind delete cluster
```

## ���� ��
```
charts/prometheus-stack/values.yaml   # Prometheus/Alertmanager ���� + Slack ������
configs/prometheus/alert-rules.yaml   # ��� �˸� ��Ģ(30s ��, 5m ����)
configs/prometheus/servicemonitor.yaml# node-exporter ��ũ������ 30s
configs/node-exporter/daemonset.yaml  # node-exporter DaemonSet
configs/node-exporter/service.yaml    # node-exporter Service
scripts/install.sh, scripts/uninstall.sh
specs/, tests/                       # ���/�׽�Ʈ ����
```

## Ʈ��������
- Slack ������ �� �ٲ� ��: �ֽ� `values.yaml`�� Helm ����� �� Alertmanager StatefulSet �����.
- ���� ��ü: Secret `slack-webhook`�� `api_url`�� �����ϸ� ����� ���� �ݿ�.
