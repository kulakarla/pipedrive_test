package config

//APIToken is the API token used for authenticating the requests to the official PipeDrive API
//Please change this to your specific API Token!
const APIToken = "863be942d8456f146e61026f7cf69dc78efda801"

//APITokenParam combines the URL parameter with the token
const APITokenParam = "?api_token=" + APIToken

//BaseURL refers to the link where the proxy API forwards its requests to
const BaseURL = "https://api.pipedrive.com/v1/deals/"
