package validate

import (
	"context"
	"encoding/json"
	"fmt"

	consumerv1 "github.com/edgefarm/provider-nats/apis/consumer/v1alpha1"
	streamsv1 "github.com/edgefarm/provider-nats/apis/stream/v1alpha1"
	natsv1 "github.com/edgefarm/provider-nats/apis/v1alpha1"

	accountv1 "github.com/edgefarm/provider-natssecrets/apis/account/v1alpha1"
	userv1 "github.com/edgefarm/provider-natssecrets/apis/user/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
)

// GENERAL

func GetUidFromDependsOnUidLabel(labels map[string]string) (string, error) {
	if val, ok := labels[dependsOnUidLabel]; ok {
		return val, nil
	}
	return "", fmt.Errorf("label %s not found", dependsOnUidLabel)
}

// CONSUMERS

func GetConsumers(ctx context.Context, client dynamic.Interface, options metav1.ListOptions) (*consumerv1.ConsumerList, error) {
	listRaw, err := client.Resource(consumerResource).List(ctx, options)
	if err != nil {
		return nil, err
	}
	list := &consumerv1.ConsumerList{}
	for _, e := range listRaw.Items {
		j, err := e.MarshalJSON()
		if err != nil {
			return nil, err
		}
		s := &consumerv1.Consumer{}
		err = json.Unmarshal(j, s)
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, *s)
	}

	return list, nil
}

func GetConsumerByName(ctx context.Context, client dynamic.Interface, name string) (*consumerv1.Consumer, error) {
	list, err := GetConsumers(ctx, client, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return nil, err
	}
	if len(list.Items) == 0 {
		return nil, ErrNotFound
	}
	if len(list.Items) > 1 {
		return nil, ErrMultipleFound
	}
	return &list.Items[0], nil
}

func GetConsumersByDependsOnUidLabel(ctx context.Context, client dynamic.Interface, uid string) (*consumerv1.ConsumerList, error) {
	return GetConsumers(ctx, client, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", dependsOnUidLabel, uid),
	})
}

func GetConsumerByUid(ctx context.Context, client dynamic.Interface, uid string) (*consumerv1.Consumer, error) {
	list, err := GetConsumers(ctx, client, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.uid=%s", uid),
	})
	if err != nil {
		return nil, err
	}
	if len(list.Items) == 0 {
		return nil, ErrNotFound
	}
	if len(list.Items) > 1 {
		return nil, ErrMultipleFound
	}
	return &list.Items[0], nil
}

// STREAMS

func GetStreams(ctx context.Context, client dynamic.Interface, options metav1.ListOptions) (*streamsv1.StreamList, error) {
	listRaw, err := client.Resource(streamResource).List(ctx, options)
	if err != nil {
		return nil, err
	}
	list := &streamsv1.StreamList{}
	for _, e := range listRaw.Items {
		j, err := e.MarshalJSON()
		if err != nil {
			return nil, err
		}
		s := &streamsv1.Stream{}
		err = json.Unmarshal(j, s)
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, *s)
	}

	return list, nil
}

func GetStreamByName(ctx context.Context, client dynamic.Interface, name string) (*streamsv1.Stream, error) {
	list, err := GetStreams(ctx, client, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return nil, err
	}
	if len(list.Items) == 0 {
		return nil, ErrNotFound
	}
	if len(list.Items) > 1 {
		return nil, ErrMultipleFound
	}
	return &list.Items[0], nil
}

func GetStreamsByDependsOnUidLabel(ctx context.Context, client dynamic.Interface, uid string) (*streamsv1.StreamList, error) {
	return GetStreams(ctx, client, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", dependsOnUidLabel, uid),
	})
}

func GetStreamByUid(ctx context.Context, client dynamic.Interface, uid string) (*streamsv1.Stream, error) {
	list, err := GetStreams(ctx, client, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.uid=%s", uid),
	})
	if err != nil {
		return nil, err
	}
	if len(list.Items) == 0 {
		return nil, ErrNotFound
	}
	if len(list.Items) > 1 {
		return nil, ErrMultipleFound
	}
	return &list.Items[0], nil
}

// USERS

func GetUsers(ctx context.Context, client dynamic.Interface, options metav1.ListOptions) (*userv1.UserList, error) {
	usersRaw, err := client.Resource(userResource).List(ctx, options)
	if err != nil {
		return nil, err
	}
	list := &userv1.UserList{}
	for _, userRaw := range usersRaw.Items {
		userJSON, err := userRaw.MarshalJSON()
		if err != nil {
			return nil, err
		}
		user := &userv1.User{}
		err = json.Unmarshal(userJSON, user)
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, *user)
	}

	return list, nil
}

func GetUserByName(ctx context.Context, client dynamic.Interface, name string) (*userv1.User, error) {
	users, err := GetUsers(ctx, client, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return nil, err
	}
	if len(users.Items) == 0 {
		return nil, ErrNotFound
	}
	if len(users.Items) > 1 {
		return nil, ErrMultipleFound
	}
	return &users.Items[0], nil
}

func GetUsersByDependsOnUidLabel(ctx context.Context, client dynamic.Interface, uid string) (*userv1.UserList, error) {
	return GetUsers(ctx, client, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", dependsOnUidLabel, uid),
	})
}

func GetUserByUid(ctx context.Context, client dynamic.Interface, uid string) (*userv1.User, error) {
	users, err := GetUsers(ctx, client, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.uid=%s", uid),
	})
	if err != nil {
		return nil, err
	}
	if len(users.Items) == 0 {
		return nil, ErrNotFound
	}
	if len(users.Items) > 1 {
		return nil, ErrMultipleFound
	}
	return &users.Items[0], nil
}

// ACCOUNTS

func GetAccounts(ctx context.Context, client dynamic.Interface, options metav1.ListOptions) (*accountv1.AccountList, error) {
	listRaw, err := client.Resource(accountResource).List(ctx, options)
	if err != nil {
		return nil, err
	}
	list := &accountv1.AccountList{}
	for _, e := range listRaw.Items {
		j, err := e.MarshalJSON()
		if err != nil {
			return nil, err
		}
		s := &accountv1.Account{}
		err = json.Unmarshal(j, s)
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, *s)
	}

	return list, nil
}

func GetAccountByName(ctx context.Context, client dynamic.Interface, name string) (*accountv1.Account, error) {
	accounts, err := GetAccounts(ctx, client, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return nil, err
	}
	if accounts.Items == nil {
		return nil, ErrNotFound
	}
	if len(accounts.Items) == 0 {
		return nil, ErrNotFound
	}
	if len(accounts.Items) > 1 {
		return nil, ErrMultipleFound
	}
	return &accounts.Items[0], nil
}

func GetAccountsByDependsOnUidLabel(ctx context.Context, client dynamic.Interface, uid string) (*accountv1.AccountList, error) {
	return GetAccounts(ctx, client, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", dependsOnUidLabel, uid),
	})
}

func GetAccountByUid(ctx context.Context, client dynamic.Interface, uid string) (*accountv1.Account, error) {
	accounts, err := GetAccounts(ctx, client, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.uid=%s", uid),
	})
	if err != nil {
		return nil, err
	}
	if len(accounts.Items) == 0 {
		return nil, ErrNotFound
	}
	if len(accounts.Items) > 1 {
		return nil, ErrMultipleFound
	}
	return &accounts.Items[0], nil
}

// PROVIDERCONFIG

func GetProviderConfigs(ctx context.Context, client dynamic.Interface, options metav1.ListOptions) (*natsv1.ProviderConfigList, error) {
	listRaw, err := client.Resource(providerConfigResource).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	list := &natsv1.ProviderConfigList{}
	for _, e := range listRaw.Items {
		j, err := e.MarshalJSON()
		if err != nil {
			return nil, err
		}
		s := &natsv1.ProviderConfig{}
		err = json.Unmarshal(j, s)
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, *s)
	}

	return list, nil
}

func GetProviderConfigByName(ctx context.Context, client dynamic.Interface, name string) (*natsv1.ProviderConfig, error) {
	providerConfigs, err := GetProviderConfigs(ctx, client, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return nil, err
	}
	if len(providerConfigs.Items) == 0 {
		return nil, ErrNotFound
	}
	if len(providerConfigs.Items) > 1 {
		return nil, ErrMultipleFound
	}
	return &providerConfigs.Items[0], nil
}

func GetProviderConfigsByDependsOnUidLabel(ctx context.Context, client dynamic.Interface, uid string) (*natsv1.ProviderConfigList, error) {
	return GetProviderConfigs(ctx, client, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", dependsOnUidLabel, uid),
	})
}

func GetProviderConfigByUid(ctx context.Context, client dynamic.Interface, uid string) (*natsv1.ProviderConfig, error) {
	providerConfigs, err := GetProviderConfigs(ctx, client, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.uid=%s", uid),
	})
	if err != nil {
		return nil, err
	}
	if len(providerConfigs.Items) == 0 {
		return nil, ErrNotFound
	}
	if len(providerConfigs.Items) > 1 {
		return nil, ErrMultipleFound
	}
	return &providerConfigs.Items[0], nil
}
