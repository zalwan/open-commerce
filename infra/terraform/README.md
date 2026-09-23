# Terraform — GCP scaffold (planned, nothing applied)

> Status: scaffold only. No GCP project, credentials, billing, or live resources exist. Nothing here has been applied; there is no drift to reconcile.

## What this provisions (when applied)

| Resource | Purpose |
|----------|---------|
| VPC + subnet + VPC connector | Private networking; Cloud Run → Cloud SQL |
| Cloud SQL Postgres 16 (private IP) | Production database (replaces local Compose `db`) |
| Artifact Registry | `api` + `web` image storage |
| Cloud Run `api` | Go backend (`DATABASE_URL` from Secret Manager) |
| Cloud Run `web` | SvelteKit storefront (`PUBLIC_API_BASE_URL` → api service) |

DB password is generated (`random_password`) and stored in Secret Manager; it never appears in git. See ADR-0004.

## Prerequisites

1. GCP project with billing: `gcloud projects create <id> --name="open-commerce"` + link billing.
2. Enable APIs: `compute`, `sqladmin`, `run`, `artifactregistry`, `secretmanager`, `servicenetworking`, `vpcaccess`.
3. Auth: `gcloud auth application-default login` (or a CI service-account key — never commit it).
4. State bucket (per project, created once): `gsutil mb -l <region> gs://oc-tfstate-<project>`.

## Build & push images first

```sh
REG=<region>-docker.pkg.dev/<project>/<env>-oc-images
gcloud auth configure-docker <region>-docker.pkg.dev
docker build -f infra/Dockerfile.api -t $REG/api:v1 .          # from repo root
docker build -f infra/Dockerfile.web -t $REG/web:v1 .          # from repo root
docker push $REG/api:v1 $REG/web:v1
```

## Plan / apply (staging example)

```sh
cd infra/terraform
terraform init \
  -backend-config="bucket=oc-tfstate-<project>" \
  -backend-config="prefix=staging"
terraform plan \
  -var=project_id=<project> \
  -var=environment=staging \
  -var=api_image=$REG/api:v1 \
  -var=web_image=$REG/web:v1
terraform apply  # same -var flags; review the plan first
```

Never commit `*.tfvars` (may hold image tags/IDs) — pass `-var` flags or use `terraform.tfvars.example` as a template (no secrets in it).

## After apply

- `terraform output web_url` → open the storefront; run `BASE=<web_url> sh tests/smoke.sh` (API at `<api_url>`).
- App auto-migrates on boot (`internal/db/migrations/`), same as Compose.
- Scale-to-zero defaults keep idle cost near Cloud SQL minimum tier only.

## Teardown

```sh
terraform destroy  # same -var flags
```

`db_deletion_protection` defaults to `true` and will block destroying the database — deliberate. Set `-var=db_deletion_protection=false` only to delete a disposable environment.

## Hardening follow-ups (not in scaffold)

Dedicated service accounts with least privilege (currently Cloud Run defaults), Cloud Armor / IAP in front of admin, budget alerts, Uptime checks on `/healthz` + `/readyz`.
