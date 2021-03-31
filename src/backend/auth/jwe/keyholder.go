package jwe

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"sync"

	jose "gopkg.in/square/go-jose.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/cuijxin/k8s-dashboard/src/backend/args"
	authApi "github.com/cuijxin/k8s-dashboard/src/backend/auth/api"
	"github.com/cuijxin/k8s-dashboard/src/backend/errors"
	syncApi "github.com/cuijxin/k8s-dashboard/src/backend/sync/api"
)

// Entries held by resource used to synchronize encryption key data.
const (
	holderMapKeyEntry  = "priv"
	holderMapCertEntry = "pub"
)

// KeyHolder is responsible for generating, storing and synchronize encryption key
// used for token generation/decryption.
type KeyHolder interface {
	// Returns encrypter instance that can be used to encrypt data.
	Encrypter() jose.Encrypter
	// Returns encryption key that can be used to decrypt data.
	Key() *rsa.PrivateKey
	// Forces refresh of encryption key synchronized with kubernetes resource (secret).
	Refresh()
}

// Implements KeyHolder interface
type rsaKeyHolder struct {
	// 256-byte random RSA key pair. Synced with a key saved in a secret.
	key          *rsa.PrivateKey
	synchronizer syncApi.Synchronizer
	mux          sync.Mutex
}

// Encrypter implements key holder interface. See KeyHolder for more information.
// Used encryption algorithms:
//  - Content encryption: AES-GCM (256)
//  - Key management: RSA-0AEP-SHA256
func (r *rsaKeyHolder) Encrypter() jose.Encrypter {
	publicKey := &r.Key().PublicKey
	encrypter, err := jose.NewEncrypter(jose.A256GCM, jose.Recipient{
		Algorithm: jose.RSA_OAEP_256,
		Key:       publicKey,
	}, nil)
	if err != nil {
		panic(err)
	}

	return encrypter
}

// Key implements key holder interface. See KeyHolder for more information.
func (r *rsaKeyHolder) Key() *rsa.PrivateKey {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.key
}

// Refresh implements key holder interface. See KeyHolder for more information.
func (r *rsaKeyHolder) Refresh() {
	r.synchronizer.Refresh()
	r.update(r.synchronizer.Get())
}

// Handler function executed by synchronizer used to store encryption key. It is
// called whenever watched object is created updated.
func (r *rsaKeyHolder) update(obj runtime.Object) {
	secret := obj.(*v1.Secret)
	priv, err := ParseRSAKey(string(secret.Data[holderMapKeyEntry]), string(secret.Data[holderMapCertEntry]))
	if err != nil {
		// Secret was probably tampered with. Update it based on local key.
		err := r.synchronizer.Update(r.getEncryptionKeyHolder())
		if err != nil {
			panic(err)
		}

		return
	}

	r.mux.Lock()
	defer r.mux.Unlock()
	r.key = priv
}

// Handler function executed by synchronizer used to store encryption key. It is
// called whenever watched object gets deleted. It is then recreated based on local key.
func (r *rsaKeyHolder) recreate(obj runtime.Object) {
	secret := obj.(*v1.Secret)
	log.Printf("Synchronized secret %s has been deleted. Recreating.", secret.Name)
	if err := r.synchronizer.Create(r.getEncryptionKeyHolder()); err != nil {
		panic(err)
	}
}

func (r *rsaKeyHolder) init() {
	r.initEncryptionKey()

	// Register event handlers
	r.synchronizer.RegisterActionHandler(r.update, watch.Added, watch.Modified)
	r.synchronizer.RegisterActionHandler(r.recreate, watch.Deleted)

	// Try to init key from synchronized object
	if obj := r.synchronizer.Get(); obj != nil {
		log.Print("Initializing JWE encryption key from synchronized object")
		r.update(obj)
		return
	}

	// Try to save generated key in a secret
	log.Printf("Storing encryption key in a secret")
	err := r.synchronizer.Create(r.getEncryptionKeyHolder())
	if err != nil && !errors.IsAlreadyExists(err) {
		panic(err)
	}
}

func (r *rsaKeyHolder) getEncryptionKeyHolder() runtime.Object {
	priv, pub := ExportRSAKeyOrDie(r.Key())
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: args.Holder.GetNamespace(),
			Name:      authApi.EncryptionKeyHolderName,
		},

		Data: map[string][]byte{
			holderMapKeyEntry:  []byte(priv),
			holderMapCertEntry: []byte(pub),
		},
	}
}

// Generates encryption key used to encrypt token payload.
func (r *rsaKeyHolder) initEncryptionKey() {
	log.Print("Generating JWE encryption key")
	r.mux.Lock()
	defer r.mux.Unlock()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	r.key = privateKey
}

// NewRSAKeyHolder creates new KeyHolder instance.
func NewRSAKeyHolder(synchronizer syncApi.Synchronizer) KeyHolder {
	holder := &rsaKeyHolder{
		synchronizer: synchronizer,
	}

	holder.init()
	return holder
}
