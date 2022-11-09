module github.com/vmware-tanzu/tanzu-framework/pinniped-components/post-deploy

go 1.18

// Right now, we depend on Kube 1.23, but Pinniped 1.20.
//
// This is because we need to depend on Pinniped v0.12.1 (because that is the
// latest Pinniped version running on TKG clusters), but Pinniped v0.12.1 does
// not contain generated Pinniped APIs/clients to go with Kube 1.23 (i.e.,
// https://github.com/vmware-tanzu/pinniped/tree/v0.12.1/generated does not
// contain a 1.23 directory).
//
// For now, we will depend on Pinniped 1.20 generated code because go mod logic
// will automatically update our dependency graph to use the later Kube version
// in the graph (i.e., our post-deploy job will use Kube v0.23.0, derived from
// the dependencies k8s.io/{api,apimachinery,client-go} v0.23.0). This seems
// like the best of the worst ways to go.
//
// In the future, we should always try to depend on the same version of Kube and
// Pinniped generated code (i.e., when we update to depend on Kube 1.24, we
// should also depend on go.pinniped.dev/generated/1.24).

require (
	github.com/go-logr/logr v1.2.2
	github.com/jetstack/cert-manager v1.1.0
	github.com/stretchr/testify v1.8.0
	github.com/vmware-tanzu/tanzu-framework/pinniped-components/tanzu-auth-controller-manager v0.0.0-00010101000000-000000000000
	go.pinniped.dev/generated/1.20/apis v0.0.0-00010101000000-000000000000
	go.pinniped.dev/generated/1.20/client v0.0.0-20220209183828-4d6a2af89419 // Commit SHA 4d6a2af89419 is tag v0.12.1.
	go.uber.org/zap v1.21.0
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/api v0.24.2
	k8s.io/apimachinery v0.24.2
	k8s.io/client-go v0.24.2
	k8s.io/klog/v2 v2.60.1
	k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9
	sigs.k8s.io/cluster-api v1.2.4
	sigs.k8s.io/controller-runtime v0.12.3
)

require (
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful v2.15.0+incompatible // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.pinniped.dev/generated/1.19/apis v0.0.0-20220310140840-61c8d5452705 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/net v0.0.0-20220617184016-355a448f1bc9 // indirect
	golang.org/x/oauth2 v0.0.0-20220608161450-d0670ef3b1eb // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220224211638-0e9765cccd65 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/apiextensions-apiserver v0.24.2 // indirect
	k8s.io/component-base v0.24.2 // indirect
	k8s.io/kube-openapi v0.0.0-20220328201542-3ee0da9b0b42 // indirect
	sigs.k8s.io/json v0.0.0-20211208200746-9f7c6b3444d2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

// Import an nested go modules have some known issues. The following replace temporarily fixes it
// https://github.com/golang/go/issues/34055
//
// Commit SHA 4d6a2af89419 is tag v0.12.1.
replace go.pinniped.dev/generated/1.20/apis => go.pinniped.dev/generated/1.19/apis v0.0.0-20220209183828-4d6a2af89419

replace github.com/vmware-tanzu/tanzu-framework/pinniped-components/tanzu-auth-controller-manager => ../tanzu-auth-controller-manager

replace sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v1.2.4