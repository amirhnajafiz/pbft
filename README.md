# todos

- [ ] Digital Certificates
- [ ] Optimal Linear PBFT
- [ ] Timers
- [ ] View-change
- [ ] Checkpointing
- [ ] Threshold Signature

```
controller -> RPC:transaction:client -> RPC:request:leader
-> client needs the controller to wait for its response
others -> RPC:reply:client -> controller
```
