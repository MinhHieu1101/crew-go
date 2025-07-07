package graphql

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"
	"net/http"
	"pkg/logger"
	"time"
	"user-service/graphql/generated"
	"user-service/internal/auth"
	cookieMw "user-service/internal/http/middleware"
	"user-service/internal/user"

	"go.uber.org/zap"
)

type Resolver struct {
	UserService user.Service
}

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, username string, email string, password string, role string) (*user.User, error) {
	return r.UserService.Register(username, email, password, role)
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*generated.Auth, error) {
	u, err := r.UserService.Authenticate(email, password)
	if err != nil {
		return nil, err
	}
	accessToken, err := auth.GenerateAccessToken(u.ID.String())
	if err != nil {
		return nil, err
	}
	refreshToken, err := auth.GenerateRefreshToken(u.ID.String())
	if err != nil {
		return nil, err
	}
	// use the cookie middleware’s FromContext to set the cookie
	w, ok := cookieMw.FromContext(ctx)
	if !ok {
		return nil, nil
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false,                // in prod use true
		SameSite: http.SameSiteLaxMode, // in prod use SameSiteStrictMode
	})
	return &generated.Auth{
		User:         u,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*user.User, error) {
	// use the auth middleware’s FromContext for claims
	claims, ok := FromContext(ctx)
	logger.Log.Debug("User query", zap.Any("claims", claims))
	if !ok || claims == nil {
		return nil, nil
	}

	return r.UserService.FindByID(claims.UserID)
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, role string) ([]*user.User, error) {
	claims, ok := FromContext(ctx)
	if !ok || claims == nil {
		return nil, nil
	}
	return r.UserService.FindByRole(role)
}

// ID is the resolver for the id field.
func (r *userResolver) ID(ctx context.Context, obj *user.User) (string, error) {
	return obj.ID.String(), nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *userResolver) CreatedAt(ctx context.Context, obj *user.User) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	type Resolver struct {
	UserService user.Service
}
*/
