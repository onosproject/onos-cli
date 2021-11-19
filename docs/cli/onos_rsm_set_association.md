## onos rsm set association

Set UE-Slice association

```
onos rsm set association [flags]
```

### Options

```
      --CuUeF1apID string    CU-UE-F1AP-ID
      --DuUeF1apID string    DU-UE-F1AP-ID
      --RanUeNgapID string   RAN-UE-NGAP-ID
      --dlSliceID string     DL Slice ID
      --drbID string         DRB-ID
      --e2NodeID string      E2 Node ID
      --eNBUeS1apID string   ENB-UE-S1AP-ID
  -h, --help                 help for association
      --no-headers           disable output headers
      --ulSliceID string     UL Slice ID
```

### Options inherited from parent commands

```
      --auth-header string       Auth header in the form 'Bearer <base64>'
      --no-tls                   if present, do not use TLS
      --service-address string   the gRPC endpoint (default "onos-rsm:5150")
      --tls-cert-path string     the path to the TLS certificate
      --tls-key-path string      the path to the TLS key
```

### SEE ALSO

* [onos rsm set](onos_rsm_set.md)	 - Set RSM resources

###### Auto generated by spf13/cobra on 25-Oct-2021