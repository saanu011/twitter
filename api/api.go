package api

import (
	r "twitter/routes"
)

var routes = r.Routes{

	r.Route{
		Name:        "UserCreate",
		Pattern:     "/users",
		Method:      "POST",
		HandlerFunc: UserCreate,
		VerifyJWT:   false,
	},

	r.Route{
		Name:        "UserLogin",
		Pattern:     "/users/login",
		Method:      "POST",
		HandlerFunc: UserLogin,
		VerifyJWT:   false,
	},

	//r.Route{
	//	Name:         "PostStoreClose",
	//	Pattern:      "/users/logout",
	//	Method:       "POST",
	//	HandlerFunc:  Logout,
	//	VerifyJWT:    true,
	//},

	r.Route{
		Name:         "PostTweet",
		Pattern:      "/users/{user_id}/tweets",
		Method:       "POST",
		HandlerFunc:  PostTweet,
		VerifyJWT:    true,
	},

	r.Route{
		Name:        "GetStoreClose",
		Pattern:     "/users/{user_id}",
		Method:      "GET",
		HandlerFunc: GetUser,
		VerifyJWT:   true,
	},

	r.Route{
		Name:        "GetStoreClose",
		Pattern:     "/users",
		Method:      "GET",
		HandlerFunc: GetUsers,
		VerifyJWT:   true,
	},

	r.Route{
		Name:         "FollowUser",
		Pattern:      "/users/{user_id}/follow",
		Method:       "POST",
		HandlerFunc:  FollowUser,
		VerifyJWT:    true,
	},

	r.Route{
		Name:         "GetUserTweets",
		Pattern:      "/users/{user_id}/tweets",
		Method:       "GET",
		HandlerFunc:  GetUserTweets,
		VerifyJWT:    true,
	},

	r.Route{
		Name:         "GetAllFollowedUsersTweets",
		Pattern:      "/users/{user_id}/tweets/followed",
		Method:       "GET",
		HandlerFunc:  GetAllFollowedUsersTweets,
		VerifyJWT:    true,
	},
}

// GetRoutes returns local variable routes which contain all methods for the API
func GetRoutes() r.Routes {
	return routes
}
