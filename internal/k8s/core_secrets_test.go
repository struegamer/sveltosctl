package k8s_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	authenticationv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	fake2 "k8s.io/client-go/kubernetes/fake"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/projectsveltos/sveltosctl/internal/k8s"
	"github.com/projectsveltos/sveltosctl/internal/logging"
)

var _ = Describe("K8S Secrets Client", func() {
	scheme := runtime.NewScheme()
	schemeErr := corev1.AddToScheme(scheme)
	_ = corev1.AddToScheme(clientsetscheme.Scheme)
	Expect(schemeErr).NotTo(HaveOccurred())
	schemeErr = authenticationv1.AddToScheme(scheme)
	_ = authenticationv1.AddToScheme(clientsetscheme.Scheme)
	Expect(schemeErr).NotTo(HaveOccurred())

	initObjects := []k8sclient.Object{}
	k8sObjects := []runtime.Object{}
	for i := 0; i < 10; i++ {
		ns := &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "acc-" + randomString(),
				Namespace: "default",
			},
		}
		initObjects = append(initObjects, ns)
		k8sObjects = append(k8sObjects, ns)
	}
	for i := 0; i < 5; i++ {
		ns := &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "acc-" + randomString(),
				Namespace: "foo",
			},
		}
		initObjects = append(initObjects, ns)
		k8sObjects = append(k8sObjects, ns)
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()
	fakeClientSet := fake2.NewClientset(k8sObjects...)

	c := &k8s.Client{}
	c.SetK8sClient(fakeClient)
	c.SetK8sClientSet(fakeClientSet)
	secretsClient := k8s.NewCoreClientWithLogger(c, logging.NewKlogTextLogger(nil))
	It("should create a secret", func() {
		err := secretsClient.CreateSecret("default", initObjects[0].GetName())
		Expect(err).NotTo(HaveOccurred())
	})
	//It("should return a token", func() {
	//	err := secretsClient.CreateServiceAccount("default", initObjects[0].GetName())
	//	Expect(err).NotTo(HaveOccurred())
	//	var token *authenticationv1.TokenRequest
	//	l, err := fakeClientSet.CoreV1().ServiceAccounts("default").List(context.TODO(), metav1.ListOptions{})
	//	fmt.Println(l)
	//	a, err := fakeClientSet.CoreV1().ServiceAccounts("default").Get(context.TODO(), initObjects[0].GetName(), metav1.GetOptions{})
	//	fmt.Println(a)
	//	token, err = secretsClient.CreateToken("default", initObjects[0].GetName(), 0)
	//	Expect(err).NotTo(HaveOccurred())
	//	Expect(token).NotTo(BeNil())
	//})
})
