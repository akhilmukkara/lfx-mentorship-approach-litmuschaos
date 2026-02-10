# Tech Notes

Technical design proposals and documentation.

## LitmusChaos Observability Design

Design proposal for adding comprehensive Prometheus metrics to LitmusChaos as part of the LFX Mentorship program (March-May 2026).

**Read the full proposal:** [DESIGN.md](DESIGN.md)

### Overview

This proposal covers:
- Complete metrics strategy for all LitmusChaos components (control plane, operator, experiments)
- Implementation architecture and approach
- Week-by-week implementation plan
- Grafana dashboard designs
- Success criteria

### Background

I'm Akhil Mukkara, a LitmusChaos contributor preparing for the "Add Prometheus Metrics to LitmusChaos Control Plane Service" mentorship.

To understand observability deeply, I built a working demo: https://github.com/akhilmukkara/prometheus-grafana-observability-demo

### Links

- **LFX Mentorship Issue:** https://github.com/litmuschaos/litmus/issues/5338
- **My Observability Demo:** https://github.com/akhilmukkara/prometheus-grafana-observability-demo
- **My Contributions:**
  - [PR #5257](https://github.com/litmuschaos/litmus/pull/5257) - UI improvement
  - [PR #4897](https://github.com/litmuschaos/litmus/pull/4897) - Helm docs

---

**Author:** Akhil Mukkara  
**Date:** February 8, 2026
