package plugins

import (
    "github.com/tietang/props/kvs"
    "github.com/tietang/zebra/meter"
    "github.com/vrischmann/go-metrics-influxdb"
    "time"
)

const (
    INFLUX_DB_INTERVAL = "metrics.export.influx.interval"
    INFLUX_DB_URL      = "metrics.export.influx.url"
    INFLUX_DB_DATABASE = "metrics.export.influx.database"
    INFLUX_DB_USERNAME = "metrics.export.influx.username"
    INFLUX_DB_PASSWORD = "metrics.export.influx.password"
)

func Start(conf kvs.ConfigSource) {
    //vrischmann/go-metrics-influxdb
    go influxdb.InfluxDB(meter.DefaultRegistry,
        conf.GetDurationDefault(INFLUX_DB_INTERVAL, time.Second*10),
        conf.GetDefault(INFLUX_DB_URL, "127.0.0.1:8086"),
        conf.GetDefault(INFLUX_DB_DATABASE, "metrics"),
        conf.GetDefault(INFLUX_DB_USERNAME, ""),
        conf.GetDefault(INFLUX_DB_PASSWORD, ""),
    )
    //yvasiyarov/go-metrics/influxdb
    //config := &influxdb.Config{
    //    Host:     conf.GetDefault(INFLUX_DB_URL, "127.0.0.1:8086"),
    //    Database: conf.GetDefault(INFLUX_DB_DATABASE, "metrics"),
    //    Username: conf.GetDefault(INFLUX_DB_USERNAME, ""),
    //    Password: conf.GetDefault(INFLUX_DB_PASSWORD, ""),
    //}
    //
    //go influxdb.Influxdb(metrics.DefaultRegistry, conf.GetDurationDefault(INFLUX_DB_INTERVAL, time.Second*10), config)
}
