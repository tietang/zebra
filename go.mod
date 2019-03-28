module github.com/tietang/zebra

go 1.12

//被墙的原因，替换golang.org源为github.com源
replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.37.2
	golang.org/x/build => github.com/golang/build v0.0.0-20190327004547-5a2224f3eb52
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190325154230-a5d413f7728c
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190321205749-f0864edee7f3
	golang.org/x/image => github.com/golang/image v0.0.0-20190321063152-3fc05d484e9f
	golang.org/x/lint => github.com/golang/lint v0.0.0-20190313153728-d0100b6bd8b3
	golang.org/x/mobile => github.com/golang/mobile v0.0.0-20190319155245-9487ef54b94a
	golang.org/x/net => github.com/golang/net v0.0.0-20190327025741-74e053c68e29
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190319182350-c85d3e98c914
	golang.org/x/perf => github.com/golang/perf v0.0.0-20190312170614-0655857e383f
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190227155943-e225da77a7e6
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190322080309-f49334f85ddc
	golang.org/x/text => github.com/golang/text v0.3.0
	golang.org/x/time => github.com/golang/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190327011446-79af862e6737
	google.golang.org/api => github.com/googleapis/google-api-go-client v0.3.0
	google.golang.org/appengine => github.com/golang/appengine v1.5.0
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20190321212433-e79c0c59cdb5
	google.golang.org/grpc => github.com/grpc/grpc-go v1.19.1
)

require (
	cloud.google.com/go v0.34.0 // indirect
	github.com/ArthurHlt/gominlog v0.0.0-20170402142412-72eebf980f46 // indirect
	github.com/Joker/hpp v1.0.0 // indirect
	github.com/Joker/jade v0.0.0-20180419144541-8828253bfc54 // indirect
	github.com/StackExchange/wmi v0.0.0-20180116203802-5d049714c4a6 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/aymerick/raymond v0.0.0-20180322193309-b565731e1464 // indirect
	github.com/codegangsta/inject v0.0.0-20140425184007-37d7f8432a3e // indirect
	github.com/deckarep/golang-set v0.0.0-20171013212420-1d4478f51bed
	github.com/denisenkom/go-mssqldb v0.0.0-20190315220205-a8ed825ac853 // indirect
	github.com/eknkc/amber v0.0.0-20171010120322-cdade1c07385 // indirect
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/fatih/color v1.4.1 // indirect
	github.com/fatih/structs v1.0.0 // indirect
	github.com/flosch/pongo2 v0.0.0-20180611110828-67f4ff8560df // indirect
	github.com/fukata/golang-stats-api-handler v1.0.0
	github.com/gavv/monotime v0.0.0-20171021193802-6f8212e8d10d // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/gin-contrib/sse v0.0.0-20170109093832-22d885f9ecc7 // indirect
	github.com/gin-gonic/gin v0.0.0-20170702092826-d459835d2b07
	github.com/go-ini/ini v1.37.0
	github.com/go-martini/martini v0.0.0-20140519164645-49411a5b6468
	github.com/go-ole/go-ole v1.2.1 // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/builder v0.2.0 // indirect
	github.com/go-xorm/core v0.6.0 // indirect
	github.com/go-xorm/xorm v0.7.0
	github.com/google/btree v1.0.0 // indirect
	github.com/google/gofuzz v0.0.0-20161122191042-44d81051d367 // indirect
	github.com/googleapis/gnostic v0.0.0-20170729233727-0c5108395e2d // indirect
	github.com/gopherjs/gopherjs v0.0.0-20180628210949-0892b62f0d9f // indirect
	github.com/hashicorp/consul v1.2.0
	github.com/hashicorp/go-cleanhttp v0.0.0-20171218145408-d5fe4b57a186 // indirect
	github.com/hashicorp/go-rootcerts v0.0.0-20160503143440-6bb64b370b90 // indirect
	github.com/hashicorp/go-uuid v1.0.1 // indirect
	github.com/hashicorp/go-version v0.0.0-20180322230233-23480c066577 // indirect
	github.com/hashicorp/memberlist v0.1.3 // indirect
	github.com/hashicorp/serf v0.0.0-20180504200640-4b67f2c2b2bb // indirect
	github.com/imkira/go-interpol v1.1.0 // indirect
	github.com/influxdata/influxdb v1.6.0 // indirect
	github.com/iris-contrib/formBinder v0.0.0-20171010160137-ad9fb86c356f // indirect
	github.com/iris-contrib/httpexpect v0.0.0-20180314041918-ebe99fcebbce // indirect
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/json-iterator/go v0.0.0-20180701071628-ab8a2e0c74be // indirect
	github.com/jtolds/gls v4.2.1+incompatible // indirect
	github.com/juju/errors v0.0.0-20170703010042-c7d06af17c68 // indirect
	github.com/juju/loggo v0.0.0-20190212223446-d976af380377 // indirect
	github.com/juju/testing v0.0.0-20180920084828-472a3e8b2073 // indirect
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/kataras/golog v0.0.0-20180321173939-03be10146386 // indirect
	github.com/kataras/iris v10.6.5+incompatible
	github.com/kataras/pio v0.0.0-20180511174041-a9733b5b6b83 // indirect
	github.com/kataras/survey v1.3.4 // indirect
	github.com/klauspost/compress v1.3.0 // indirect
	github.com/klauspost/cpuid v0.0.0-20170728055534-ae7887de9fa5 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/lafikl/consistent v0.0.0-20171026144656-ea75672e5603
	github.com/lestrrat/go-envload v0.0.0-20180220120943-6ed08b54a570 // indirect
	github.com/lestrrat/go-file-rotatelogs v0.0.0-20180223000712-d3151e2a480f
	github.com/lestrrat/go-strftime v0.0.0-20180220042222-ba3bf9c1d042 // indirect
	github.com/lib/pq v1.0.0 // indirect
	github.com/mattn/go-colorable v0.0.0-20170210172801-5411d3eea597
	github.com/mattn/go-isatty v0.0.0-20170307163044-57fdcb988a5c // indirect
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/microcosm-cc/bluemonday v0.0.0-20180621201946-f0761eb8ed07 // indirect
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1 // indirect
	github.com/mitchellh/go-homedir v0.0.0-20180523094522-3864e76763d9 // indirect
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/mitchellh/mapstructure v0.0.0-20171017171808-06020f85339e // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pascaldekloe/goe v0.1.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20181016184325-3113b8401b8a
	github.com/rifflock/lfshook v0.0.0-20180227222202-bf539943797a
	github.com/ryanuber/columnize v0.0.0-20170703205827-abc90934186a // indirect
	github.com/samuel/go-zookeeper v0.0.0-20180130194729-c4fab1ac1bec
	github.com/satori/go.uuid v0.0.0-20180103174451-36e9d2ebbde5 // indirect
	github.com/sergi/go-diff v1.0.0 // indirect
	github.com/shirou/gopsutil v0.0.0-20180625081143-4a180b209f5f
	github.com/shirou/w32 v0.0.0-20160930032740-bb4de0191aa4 // indirect
	github.com/shurcooL/sanitized_anchor_name v0.0.0-20170918181015-86672fcb3f95 // indirect
	github.com/sirupsen/logrus v1.2.0
	github.com/smartystreets/assertions v0.0.0-20180301161246-7678a5452ebe // indirect
	github.com/smartystreets/goconvey v0.0.0-20170602164621-9e8dc3f972df
	github.com/smartystreets/gunit v0.0.0-20180314194857-6f0d6275bdcd // indirect
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/stretchr/testify v1.3.0 // indirect
	github.com/tebeka/strftime v0.0.0-20140926081919-3f9c7761e312 // indirect
	github.com/thoas/stats v0.0.0-20170926101542-37829025d224
	github.com/tietang/assert v0.0.0-20160910015056-6961f642d923
	github.com/tietang/go-eureka-client v0.0.0-20171116042000-3f12f7db4199
	github.com/tietang/go-utils v0.0.0-20180420232328-76758c2288ca
	github.com/tietang/godaemon v0.0.0-20160320101618-2d183393d9ee
	github.com/tietang/hystrix-go v0.0.0-20170922014527-a984df1911a5
	github.com/tietang/props v2.2.0+incompatible
	github.com/tietang/stats v0.0.0-20171114031414-9f32ebcae985
	github.com/ugorji/go v0.0.0-20170215201144-c88ee250d022 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v0.0.0-20170224212429-dcecefd839c4 // indirect
	github.com/vrischmann/go-metrics-influxdb v0.0.0-20160917065939-43af8332c303
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.1.0 // indirect
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0 // indirect
	github.com/yudai/gojsondiff v1.0.0 // indirect
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82 // indirect
	github.com/yudai/pp v2.0.1+incompatible // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
	golang.org/x/net v0.0.0-20190327091125-710a502c58a2
	gopkg.in/AlecAivazis/survey.v1 v1.5.3 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.1 // indirect
	gopkg.in/ini.v1 v1.37.0
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.0-20180526075726-670777b536d3 // indirect
	gopkg.in/square/go-jose.v2 v2.1.6
	k8s.io/api v0.0.0-20180308224125-73d903622b73 // indirect
	k8s.io/apimachinery v0.0.0-20180707232508-bce280dade67
	k8s.io/client-go v7.0.0+incompatible
)
