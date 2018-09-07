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





type Movies struct {
	Page         int `json:"page"`
	TotalResults int `json:"total_results"`
	TotalPages   int `json:"total_pages"`
	Results      []MovieResult `json:"results"`
}



type MovieResult struct{
	VoteCount        int     `json:"vote_count"`
	ID               int     `json:"id"`
	Video            bool    `json:"video"`
	VoteAverage      float64 `json:"vote_average"`
	Title            string  `json:"title"`
	Popularity       float64 `json:"popularity"`
	PosterPath       string  `json:"poster_path"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	GenreIds         []int   `json:"genre_ids"`
	BackdropPath     string  `json:"backdrop_path"`
	Adult            bool    `json:"adult"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
}


type MongoDBResponse struct {
	ID         int    `json:"id" bson:"id"`
	Title      string `json:"title" bson:"title"`
	Language   string `json:"language" bson:"language"`
	Publisher  string `json:"publisher" bson:"publisher"`
	Authorname string `json:"authorname" bson:"authorname"`
	Genre      string `json:"genre" bson:"genre"`
	Year       string `json:"year" bson:"year"`
	URL        string `json:"url" bson:"url"`
}

type LibraryRequest struct {
	OriginalQuery string `json:"originalQuery"`
	Action        string `json:"action"`
	Intent        string `json:"intent"`
	Entities      map[string]string `json:"entities"`
	Response      struct {
			      Text   string `json:"text"`
			      Speech string `json:"speech"`
		      } `json:"response"`
	IsResolved bool   `json:"isResolved"`
	SessionID  string `json:"sessionId"`
}