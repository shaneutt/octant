package objectstore

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/internal/cluster/fake"
	"github.com/vmware/octant/internal/gvk"
	"github.com/vmware/octant/internal/testutil"
)

func Test_factoriesCache(t *testing.T) {
	const namespaceName = "test"

	controller := gomock.NewController(t)
	defer controller.Finish()

	dynamicClient := fake.NewMockDynamicInterface(controller)

	client := fake.NewMockClientInterface(controller)
	client.EXPECT().
		DynamicClient().
		Return(dynamicClient, nil)

	c := initFactoriesCache()

	ctx := context.Background()
	factory, err := initDynamicSharedInformerFactory(ctx, client, namespaceName)
	require.NoError(t, err)

	c.set(namespaceName, factory)

	got, isFound := c.get(namespaceName)
	require.True(t, isFound)
	require.Equal(t, factory, got)

	c.delete(namespaceName)
	_, isFound = c.get(namespaceName)
	require.False(t, isFound)
}

func Test_accessCache(t *testing.T) {
	c := initAccessCache()

	key := accessKey{
		Namespace: "test",
		Group:     "group",
		Resource:  "resource",
		Verb:      "list",
	}

	c.set(key, true)

	got, isFound := c.get(key)
	require.True(t, isFound)
	require.True(t, got)
}

func Test_seenGVKsCache(t *testing.T) {
	c := initSeenGVKsCache()
	c.setSeen("test", gvk.Pod, true)

	tests := []struct {
		name      string
		namespace string
		gvk       schema.GroupVersionKind
		expected  bool
	}{
		{
			name:      "gvk that has been seen",
			namespace: "test",
			gvk:       gvk.Pod,
			expected:  true,
		},
		{
			name:      "namespace that has not been seen",
			namespace: "other",
			gvk:       gvk.Pod,
			expected:  false,
		},
		{
			name:      "gvk that has not been seen",
			namespace: "test",
			gvk:       gvk.Deployment,
			expected:  false,
		},
	}

	for i := range tests {
		test := tests[i]

		t.Run(test.name, func(t *testing.T) {
			got := c.hasSeen(test.namespace, test.gvk)
			require.Equal(t, test.expected, got)
		})
	}
}

func Test_cachedObjectsCache(t *testing.T) {
	c := initCachedObjectsCache()

	pod := testutil.ToUnstructured(t, testutil.CreatePod("pod"))

	c.update("test", gvk.Pod, pod)

	items := c.list("test", gvk.Pod)
	require.Equal(t, []*unstructured.Unstructured{pod}, items)

	items = c.list("test", gvk.Deployment)
	require.Empty(t, items)

	items = c.list("other", gvk.Pod)
	require.Empty(t, items)

	c.delete("test", gvk.Pod, pod)

	items = c.list("test", gvk.Pod)
	require.Empty(t, items)
}

func Test_watchedGVKsCache(t *testing.T) {
	c := initWatchedGVKsCache()

	c.setWatched("test", gvk.Pod)
	assert.True(t, c.isWatched("test", gvk.Pod))
	assert.False(t, c.isWatched("test", gvk.Deployment))
	assert.False(t, c.isWatched("other", gvk.Pod))

}
