package tests

import (
	sso "github.com/Rasikrr/protobuff/protos/gen/go/sso"
	"github.com/brianvoe/gofakeit"
	"github.com/go-playground/assert/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"sso/tests/suite"
	"strconv"
	"testing"
)

const (
	emptyAppId     = 0
	appId          = 3
	appSecret      = "test-secret"
	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := generateRandomPassword()

	respReg, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEqual(t, respReg.GetUserId(), 0)

	respLogin, err := st.AuthClient.Login(ctx, &sso.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appId,
	})

	require.NoError(t, err)

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)

	assert.Equal(t, ok, true)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appId, int(claims["app_id"].(float64)))

	assert.Equal(t, tokenParsed.Valid, true)

}

func TestRegisterLogin_DuplicateRegister(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := generateRandomPassword()

	respReg, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEqual(t, respReg.GetUserId(), 0)

	_, err = st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid credentials")
}

func TestRegisterLogin_Register_InvalidCredentials(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		num         int
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			num:         1,
			name:        gofakeit.Name(),
			email:       "",
			password:    generateRandomPassword(),
			expectedErr: "invalid email",
		},
		{
			num:         2,
			name:        gofakeit.Name(),
			email:       gofakeit.Email(),
			password:    "1",
			expectedErr: "password must contain at least one upper, one digit and one spec.symbol",
		},
		{
			num:         3,
			name:        gofakeit.Name(),
			email:       "",
			password:    "",
			expectedErr: "invalid email",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name+" "+strconv.Itoa(tt.num), func(t *testing.T) {
			t.Parallel()
			_, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestRegisterLogin_Login_InvalidCredentials(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		num         int
		name        string
		email       string
		password    string
		appId       int32
		expectedErr string
	}{
		{
			num:         1,
			name:        gofakeit.Name(),
			email:       gofakeit.Email(),
			password:    "",
			appId:       appId,
			expectedErr: "invalid credentials",
		},
		{
			num:         2,
			name:        gofakeit.Name(),
			email:       "",
			password:    generateRandomPassword(),
			appId:       appId,
			expectedErr: "invalid email",
		},
		{
			num:         3,
			name:        gofakeit.Name(),
			email:       "",
			password:    "",
			appId:       appId,
			expectedErr: "invalid email",
		},
		{
			num:         4,
			name:        gofakeit.Name(),
			email:       gofakeit.Email(),
			password:    generateRandomPassword(),
			appId:       emptyAppId,
			expectedErr: "app_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" "+strconv.Itoa(tt.num), func(t *testing.T) {
			t.Parallel()
			_, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
				Email:    gofakeit.Email(),
				Password: generateRandomPassword(),
			})
			require.NoError(t, err)
			_, err = st.AuthClient.Login(ctx, &sso.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appId,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}

}

func generateRandomPassword() string {
	password := gofakeit.Password(
		true,
		true,
		true,
		true,
		false,
		passDefaultLen,
	)
	return password
}
