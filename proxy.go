package main

type UserCluster struct {
	Username    string
	ClusterName string
	ApiURL      string
	TokenHash   string
}

type UserClusters struct {
	cacheByToken       map[string]*UserCluster // by token hash
	cacheByUserCluster map[string]*UserCluster // by userCluster hash
}

func (c *UserClusters) Url(token string) (string, error) {
	userCluster, ok := c.cacheByToken[tokenHash(token)]
	if ok {
		return userCluster.ApiURL, nil
	}
	userCluster, err := c.loadCluster(token)
	if err != nil {
		return "", err
	}
	return userCluster.ApiURL, nil
}

func (c *UserClusters) loadCluster(token string) (*UserCluster, error) {
	//TODO load all member clusters

	//TODO iterate member clusters and check whoami, so we know what cluster API URL and Username
	username := ""
	apiUrl := ""
	clusterName := ""

	// Cleanup existing cached tokens user clusters if any
	ucHash := userClusterHash(username, apiUrl)
	userCluster, ok := c.cacheByUserCluster[ucHash]
	if ok {
		c.cacheByToken[userCluster.TokenHash] = nil
		c.cacheByUserCluster[ucHash] = nil
	}

	// Create UserCluster and put to cache
	userCluster = &UserCluster{
		Username:    username,
		ClusterName: clusterName,
		ApiURL:      apiUrl,
		TokenHash:   tokenHash(token),
	}
	c.cacheByToken[userCluster.TokenHash] = userCluster
	c.cacheByUserCluster[ucHash] = userCluster

	return userCluster, nil
}

func tokenHash(token string) string {
	//TODO
	return token
}

func userClusterHash(username, apiURL string) string {
	//TODO
	return username + apiURL
}
