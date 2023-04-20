package cloudflare

import (
	"io"
	"log"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"k8s.io/kops/dnsprovider/pkg/dnsprovider"
	"k8s.io/kops/dnsprovider/pkg/dnsprovider/rrstype"
)

var _ dnsprovider.Interface = Interface{}

const (
	// cloudFlareCreate is a ChangeAction enum value
	CloudFlareCreate = "CREATE"
	// cloudFlareDelete is a ChangeAction enum value
	CloudFlareDelete = "DELETE"
	// cloudFlareUpdate is a ChangeAction enum value
	CloudFlareUpdate = "UPDATE"
	// defaultCloudFlareRecordTTL 1 = automatic
	DefaultCloudFlareRecordTTL = 1
	// Provider Name for Provider init
	ProviderName = "cloudflare"
)

func init() {
	dnsprovider.RegisterDNSProvider(ProviderName, func(config io.Reader) (dnsprovider.Interface, error) {
		client, err := newApiClient()
		if err != nil {
			return nil, err
		}
		return NewProvider(client), nil
	})
}

// DNS implements dnsprovider.Interface
type Interface struct {
	client *cloudflare.API
}

// NewProvider returns an implementation of dnsprovider.Interface
func NewProvider(client *cloudflare.API) dnsprovider.Interface {
	return &Interface{client: client}
}

// Zones returns an implementation of dnsprovider.Zones
func (d Interface) Zones() (dnsprovider.Zones, bool) {
	return &zones{
		client: d.client,
	}, true
}

// zones is an implementation of dnsprovider.Zones
type zones struct {
	client *cloudflare.API
}

func newApiClient() (*cloudflare.API, error) {
	// Require a scoped API token
	api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	return api, nil

}

// Add adds a new DNS zone
func (z *zones) Add(newZone dnsprovider.Zone) (dnsprovider.Zone, error) {
	return nil, nil
}

// List returns a list of all dns zones
func (z *zones) List() ([]dnsprovider.Zone, error) {
	return nil, nil
}

// Remove deletes a zone
func (z *zones) Remove(zone dnsprovider.Zone) error {
	return nil
}

// New returns a new implementation of dnsprovider.Zone
func (z *zones) New(name string) (dnsprovider.Zone, error) {
	return nil, nil
}

// zone implements dnsprovider.Zone
type zone struct {
	name   string
	client *cloudflare.API
}

// Name returns the Name of a dns zone
func (z *zone) Name() string {
	return z.name
}

// ID returns the name of a dns zone, in DO the ID is the name
func (z *zone) ID() string {
	return z.name
}

// ResourceRecordSet returns an implementation of dnsprovider.ResourceRecordSets
func (z *zone) ResourceRecordSets() (dnsprovider.ResourceRecordSets, bool) {
	return nil, false
}

// resourceRecordSets implements dnsprovider.ResourceRecordSet
type resourceRecordSets struct {
	zone   *zone
	client *cloudflare.API
}

// List returns a list of dnsprovider.ResourceRecordSet
func (r *resourceRecordSets) List() ([]dnsprovider.ResourceRecordSet, error) {
	return nil, nil
}

// Get returns a list of dnsprovider.ResourceRecordSet that matches the name parameter
func (r *resourceRecordSets) Get(name string) ([]dnsprovider.ResourceRecordSet, error) {
	return nil, nil
}

// New returns an implementation of dnsprovider.ResourceRecordSet
func (r *resourceRecordSets) New(name string, rrdatas []string, ttl int64, rrstype rrstype.RrsType) dnsprovider.ResourceRecordSet {
	return nil
}

// StartChangeset returns an implementation of dnsprovider.ResourceRecordChangeset
func (r *resourceRecordSets) StartChangeset() dnsprovider.ResourceRecordChangeset {
	return nil
}

// Zone returns the associated implementation of dnsprovider.Zone
func (r *resourceRecordSets) Zone() dnsprovider.Zone {
	return r.zone
}

// recordRecordSet implements dnsprovider.ResourceRecordSet which represents
// a single record associated with a zone
type resourceRecordSet struct {
	name       string
	data       []string
	ttl        int
	recordType rrstype.RrsType
}

// Name returns the name of a resource record set
func (r *resourceRecordSet) Name() string {
	return r.name
}

// Rrdatas returns a list of data associated with a resource record set
// in DO this is almost always the IP of a record
func (r *resourceRecordSet) Rrdatas() []string {
	return r.data
}

// Ttl returns the time-to-live of a record
func (r *resourceRecordSet) Ttl() int64 {
	return int64(r.ttl)
}

// Type returns the type of record a resource record set is
func (r *resourceRecordSet) Type() rrstype.RrsType {
	return r.recordType
}
