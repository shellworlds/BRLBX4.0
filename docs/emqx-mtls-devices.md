# EMQX: TLS and client certificates for kitchen devices

## Goal

Kitchen controllers connect to MQTT over **TLS** with **mutual TLS**: the broker verifies a client certificate issued by your device CA.

## Server side (EMQX)

1. Terminate TLS on the MQTT listener (port 8883 typical) with a server certificate trusted by devices (public CA or private PKI distributed to devices).
2. Enable **peer verification** and configure the **CA** that signed device client certificates (same CA as `IOT_DEVICE_CA_CERT_FILE` in `iot-ingestion`).
3. Map client certificate CN or username to ACLs so each device only publishes to `borelsigma/kitchen/{kitchen_id}/telemetry`.

## Device onboarding

1. Device generates a private key and **CSR** (PEM).
2. Operator or automation calls `POST /api/v1/internal/devices/register` on `iot-ingestion` with header `X-Internal-Token: <INTERNAL_DEVICE_TOKEN>` and body:

```json
{
  "kitchen_id": "<uuid>",
  "label": "controller-1",
  "csr_pem": "-----BEGIN CERTIFICATE REQUEST-----\n..."
}
```

3. When `IOT_DEVICE_CA_CERT_FILE` and `IOT_DEVICE_CA_KEY_FILE` are mounted, the service returns a **client certificate PEM** (`status: active`). Otherwise the CSR is stored with `status: pending_signing`.

4. Install the issued cert + device private key on the controller; configure MQTT client to use them for TLS.

## Operations

- Rotate the device CA before certificates expire; re-enroll devices with new CSRs.
- Keep `INTERNAL_DEVICE_TOKEN` in Sealed Secrets; restrict network access to the registration endpoint.
