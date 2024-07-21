package third_party

import (
	"os"

	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"

	"core/database"
	"core/models"
	"core/repositories"
)

type Supertoken struct {
}

func (s Supertoken) Initialize() {
	apiBasePath := "/auth"
	websiteBasePath := "/auth"
	antiCsrf := session.AntiCSRF_VIA_CUSTOM_HEADER
	cookieDomain := ".localhost"
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: os.Getenv("SUPERTOKENS_CONNECTION_URI"),
			APIKey:        os.Getenv("SUPERTOKENS_API_KEY"),
		},
		AppInfo: supertokens.AppInfo{
			AppName:         os.Getenv("APP_NAME"),
			APIDomain:       os.Getenv("SUPERTOKENS_API_DOMAIN"),
			WebsiteDomain:   os.Getenv("SUPERTOKENS_WEBSITE_DOMAIN"),
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			dashboard.Init(nil),
			session.Init(&sessmodels.TypeInput{
				ExposeAccessTokenToFrontendInCookieBasedAuth: true,
				AntiCsrf:     &antiCsrf,
				CookieDomain: &cookieDomain,
			}),
			emailpassword.Init(&epmodels.TypeInput{
				Override: &epmodels.OverrideStruct{
					Functions: func(originalImplementation epmodels.RecipeInterface) epmodels.RecipeInterface {
						originalSignUp := *originalImplementation.SignUp
						(*originalImplementation.SignUp) = func(email, password, tenantId string, userContext supertokens.UserContext) (epmodels.SignUpResponse, error) {
							response, err := originalSignUp(email, password, tenantId, userContext)
							if err != nil {
								return epmodels.SignUpResponse{}, err
							}
							if response.OK != nil {
								if response.EmailAlreadyExistsError != nil {
									return response, err
								}
								user := response.OK.User
								repo := repositories.NewTenantRepositoryFromDB(database.GetCentralConnection())
								repo.Create(models.Tenant{
									DisplayName: &user.Email,
									AuthID:      user.ID,
								})
								return response, err

							}
							return response, nil
						}
						return originalImplementation
					},
				},
			}),
		},
	})

	if err != nil {
		panic(err.Error())
	}
}
