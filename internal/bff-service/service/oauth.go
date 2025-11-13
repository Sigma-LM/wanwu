package service

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	iam_service "github.com/UnicomAI/wanwu/api/proto/iam-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	oauth2_util "github.com/UnicomAI/wanwu/internal/bff-service/pkg/oauth2-util"
	jwt_util "github.com/UnicomAI/wanwu/pkg/jwt-util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func OAuthAuthorize(ctx *gin.Context, req *request.OAuthRequest, userID string) (string, string, error) {
	oauthApp, err := iam.GetOauthApp(ctx, &iam_service.GetOauthAppReq{
		ClientId: req.ClientID,
	})
	if err != nil {
		return "", "", err
	}
	if !oauthApp.Status {
		return "", "", fmt.Errorf("client id: %v has been disabled", oauthApp.ClientId)
	}
	if oauthApp.ClientId != req.ClientID || (req.RedirectURI != "" && req.RedirectURI != oauthApp.RedirectUri) {
		return "", "", fmt.Errorf("client id: %v or redirecturi: %v missmatch", req.ClientID, req.RedirectURI)
	}
	//code save to redis
	code := uuid.NewString()
	if err := oauth2_util.SaveCode(ctx, code, oauth2_util.CodePayload{
		ClientID: req.ClientID,
		UserID:   userID,
	}); err != nil {
		return "", "", fmt.Errorf("%v get auth info err:%v", req.ClientID, err)
	}
	return oauthApp.RedirectUri, code, nil
}

func OAuthToken(ctx *gin.Context, req *request.OAuthTokenRequest) (*response.OAuthTokenResponse, error) {
	codePayload, err := oauth2_util.ValidateCode(ctx, req.Code, req.ClientID)
	if err != nil {
		return nil, fmt.Errorf("validate code timeout err: %v", err)
	}
	oauthApp, err := iam.GetOauthApp(ctx, &iam_service.GetOauthAppReq{
		ClientId: req.ClientID,
	})
	if err != nil {
		return nil, fmt.Errorf("%v get auth info err: %v", req.ClientID, err)
	}
	err = oauthValidateCode(req.ClientID, req.ClientSecret, req.RedirectURI, codePayload, oauthApp)
	if err != nil {
		return nil, err
	}
	user, err := iam.GetUserInfo(ctx, &iam_service.GetUserInfoReq{
		UserId: codePayload.UserID,
		OrgId:  "",
	})
	if err != nil {
		return nil, err
	}
	//access token
	scopes := []string{} //预留scope处理
	accessToken, err := oauth2_util.GenerateAccessToken(user.UserId, req.ClientID, scopes, oauth2_util.AccessTokenTimeout)
	if err != nil {
		return nil, err
	}
	//id token
	idToken, err := oauth2_util.GenerateIDToken(user.UserId, user.UserName, req.ClientID, oauth2_util.IDTokenTimeout)
	if err != nil {
		return nil, err
	}
	//refresh token
	refreshToken, err := oauth2_util.GenerateRefreshToken(ctx, user.UserId, req.ClientID, oauth2_util.RefreshTokenExpiration)
	if err != nil {
		return nil, err
	}
	return &response.OAuthTokenResponse{
		AccessToken:  accessToken,
		ExpiresIn:    oauth2_util.AccessTokenTimeout,
		TokenType:    "Bearer",
		IDToken:      idToken,
		RefreshToken: refreshToken,
		Scope:        scopes,
	}, nil
}

func OAuthRefresh(ctx *gin.Context, req *request.OAuthRefreshRequest) (*response.OAuthRefreshTokenResponse, error) {
	refreshPayload, err := oauth2_util.ValidateRefreshToken(ctx, req.RefreshToken, req.ClientID)
	if err != nil {
		return nil, err
	}
	oauthApp, err := iam.GetOauthApp(ctx, &iam_service.GetOauthAppReq{
		ClientId: req.ClientID,
	})
	if err != nil {
		return nil, err
	}
	if !oauthApp.Status {
		return nil, fmt.Errorf("client id: %v has been disabled", oauthApp.ClientId)
	}
	if req.ClientSecret != oauthApp.ClientSecret {
		return nil, fmt.Errorf("clinetId:%v or clientSecret missmatch", req.ClientID)
	}
	scopes := []string{} //scopes处理预留
	//new access token
	accessToken, err := oauth2_util.GenerateAccessToken(refreshPayload.UserID, req.ClientID, scopes, oauth2_util.AccessTokenTimeout)
	if err != nil {
		return nil, err
	}
	//new refresh token
	refreshToken, err := oauth2_util.GenerateRefreshToken(ctx, refreshPayload.UserID, refreshPayload.ClientID, oauth2_util.RefreshTokenExpiration)
	if err != nil {
		return nil, err
	}
	return &response.OAuthRefreshTokenResponse{
		AccessToken:  accessToken,
		ExpiresAt:    strconv.Itoa(int(time.Now().Add(time.Duration(jwt_util.UserTokenTimeout) * time.Second).UnixMilli())),
		RefreshToken: refreshToken,
	}, nil
}

func OAuthConfig(ctx *gin.Context) (*response.OAuthConfig, error) {
	issuer, err := oauth2_util.GetIssuer()
	if err != nil {
		return nil, err
	}
	return &response.OAuthConfig{
		Issuer:           issuer,
		AuthEndpoint:     issuer + "/user/api/v1" + "/oauth/code/authorize",
		TokenEndpoint:    issuer + "/user/api/openapi/v1" + "/oauth/code/token",
		JwksUri:          issuer + "/user/api/openapi/v1" + "/oauth/jwks",
		UserInfoEndpoint: issuer + "/user/api/openapi/v1" + "/oauth/userinfo",
		ResponseTypes:    []string{"code"},
		IDtokenSignAlg:   []string{"RS256"},
		SubjectTypes:     []string{"public"},
	}, nil
}

func OAuthJWKS(ctx *gin.Context) (*response.OAuthJWKS, error) {
	jwk, err := oauth2_util.GetJWK()
	if err != nil {
		return nil, err
	}
	return &response.OAuthJWKS{Keys: []oauth2_util.JWK{jwk}}, nil
}

func OAuthGetUserInfo(ctx *gin.Context, userID string) (*response.OAuthGetUserInfo, error) {
	user, err := iam.GetUserInfo(ctx, &iam_service.GetUserInfoReq{
		UserId: userID,
		OrgId:  "",
	})
	if err != nil {
		return nil, err
	}
	issuer, err := oauth2_util.GetIssuer()
	if err != nil {
		return nil, err
	}
	avatar := cacheUserAvatar(ctx, user.AvatarPath)
	avatarUri, err := url.JoinPath(issuer, "/user/api", avatar.Path)
	if err != nil {
		return nil, err
	}
	return &response.OAuthGetUserInfo{
		UserID:    user.UserId,
		Username:  user.UserName,
		Email:     user.Email,
		Nickname:  user.NickName,
		Phone:     user.Phone,
		Gender:    user.Gender,
		AvatarUri: avatarUri,
		Remark:    user.Remark,
		Company:   user.Company,
	}, nil
}

func oauthValidateCode(clientID, clientSecret, redirectUri string, codePayload oauth2_util.CodePayload, appInfo *iam_service.OauthApp) error {
	if !appInfo.Status {
		return fmt.Errorf("client id: %v has been disabled", codePayload.ClientID)
	}
	if codePayload.ClientID != clientID { //两次传的不一样
		return fmt.Errorf("client_id mismatch: expected %v, got %v", codePayload.ClientID, clientID)
	}

	if appInfo.ClientSecret != clientSecret {
		return fmt.Errorf("client_secret error for client_id: %v", codePayload.ClientID)
	}

	if redirectUri != "" && redirectUri != appInfo.RedirectUri {
		return fmt.Errorf("redirect_uri err: got %v", redirectUri)
	}
	return nil
}

func CreateOauthApp(ctx *gin.Context, userId string, req *request.CreateOauthAppReq) error {
	_, err := iam.CreateOauthApp(ctx, &iam_service.CreateOauthAppReq{
		UserId:      userId,
		Name:        req.Name,
		Desc:        req.Desc,
		RedirectUri: req.RedirectURI,
	})
	if err != nil {
		return err
	}
	return nil
}

func DeleteOauthApp(ctx *gin.Context, req *request.DeleteOauthAppReq) error {
	_, err := iam.DeleteOauthApp(ctx, &iam_service.DeleteOauthAppReq{
		ClientId: req.ClientID,
	})
	if err != nil {
		return err
	}
	return nil
}

func UpdateOauthApp(ctx *gin.Context, req *request.UpdateOauthAppReq) error {
	_, err := iam.UpdateOauthApp(ctx, &iam_service.UpdateOauthAppReq{
		ClientId:    req.ClientID,
		Name:        req.Name,
		Desc:        req.Desc,
		RedirectUri: req.RedirectURI,
	})
	if err != nil {
		return err
	}
	return nil
}

func GetOauthAppList(ctx *gin.Context, userId string, pageNo, pageSize int32) (*response.PageResult, error) {
	resp, err := iam.GetOauthAppList(ctx, &iam_service.GetOauthAppListReq{
		UserId:   userId,
		PageNo:   pageNo,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}
	var retList []*response.OAuthAppInfo
	for _, app := range resp.Apps {
		retList = append(retList, &response.OAuthAppInfo{
			ClientID:     app.ClientId,
			Name:         app.Name,
			Desc:         app.Desc,
			ClientSecret: app.ClientSecret,
			RedirectURI:  app.RedirectUri,
			Status:       app.Status,
		})
	}
	return &response.PageResult{
		List:     retList,
		Total:    resp.Total,
		PageNo:   int(pageNo),
		PageSize: int(pageSize),
	}, nil
}

func UpdateOauthAppStatus(ctx *gin.Context, req *request.UpdateOauthAppStatusReq) error {
	_, err := iam.UpdateOauthAppStatus(ctx, &iam_service.UpdateOauthAppStatusReq{
		ClientId: req.ClientID,
		Status:   req.Status,
	})
	if err != nil {
		return err
	}
	return nil
}
