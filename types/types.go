package types

type DialogFlowRequest struct {
	ResponseID  string `json:"responseId"`
	QueryResult struct {
		QueryText  string `json:"queryText"`
		Action     string `json:"action"`
		Parameters struct {
			Author   string `json:"Author"`
			Category string `json:"Category"`
		} `json:"parameters"`
		AllRequiredParamsPresent bool `json:"allRequiredParamsPresent"`
		FulfillmentMessages      []struct {
			Text struct {
				Text []string `json:"text"`
			} `json:"text"`
		} `json:"fulfillmentMessages"`
		Intent struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
		} `json:"intent"`
		IntentDetectionConfidence int    `json:"intentDetectionConfidence"`
		LanguageCode              string `json:"languageCode"`
	} `json:"queryResult"`
}

type DialogFlowResponse struct {
	FulfillmentText     string `json:"fulfillmentText"`
	FulfillmentMessages []struct {
		Card struct {
			     Title    string `json:"title"`
			     Subtitle string `json:"subtitle"`
			     ImageURI string `json:"imageUri"`
			     Buttons  []struct {
				     Text     string `json:"text"`
				     Postback string `json:"postback"`
			     } `json:"buttons"`
		     } `json:"card"`
	} `json:"fulfillmentMessages"`
	Source  string `json:"source"`
	Payload struct {
				    Google struct {
						   ExpectUserResponse bool `json:"expectUserResponse"`
						   RichResponse       struct {
									      Items []struct {
										      SimpleResponse struct {
													     TextToSpeech string `json:"textToSpeech"`
												     } `json:"simpleResponse"`
									      } `json:"items"`
								      } `json:"richResponse"`
					   } `json:"google"`
				    Facebook struct {
						   Text string `json:"text"`
					   } `json:"facebook"`
				    Slack struct {
						   Text string `json:"text"`
					   } `json:"slack"`
			    } `json:"payload"`
	OutputContexts []struct {
		Name          string `json:"name"`
		LifespanCount int    `json:"lifespanCount"`
		Parameters    struct {
				      Param string `json:"param"`
			      } `json:"parameters"`
	} `json:"outputContexts"`
	FollowupEventInput struct {
				    Name         string `json:"name"`
				    LanguageCode string `json:"languageCode"`
				    Parameters   struct {
							 Param string `json:"param"`
						 } `json:"parameters"`
			    } `json:"followupEventInput"`
}