package discovery

type Instance struct {
    InstanceId   string
    Name         string
    AppName      string
    AppGroupName string
    //
    //
    Tags             []string
    Params           map[string]string
    InstanceType     string
    ExternalInstance interface{}
    //
    Scheme               string
    Address              string
    Port                 string
    HealthCheckUrl       string
    Status               string
    LastUpdatedTimestamp int
    OverriddenStatus     string
}
type Service struct {
    Name      string
    Labels    map[string]string
    Instances []*Instance
}
