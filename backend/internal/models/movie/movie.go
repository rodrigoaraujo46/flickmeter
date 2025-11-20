package movie

type Genre struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type ProductionCompany struct {
	Id            int32  `json:"id"`
	LogoPath      string `json:"logo_path"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

type ProductionCountry struct {
	ISO31661 string `json:"iso_3166_1"`
	Name     string `json:"name"`
}

type SpokenLanguage struct {
	EnglishName string `json:"english_name"`
	ISO6391     string `json:"iso_639_1"`
	Name        string `json:"name"`
}

type BelongsToCollection struct {
	Id           int32  `json:"id"`
	Name         string `json:"name"`
	PosterPath   string `json:"poster_path"`
	BackdropPath string `json:"backdrop_path"`
}

type Movie struct {
	Adult               bool                 `json:"adult"`
	BackdropPath        string               `json:"backdrop_path"`
	BelongsToCollection *BelongsToCollection `json:"belongs_to_collection,omitempty"`
	Budget              int32                `json:"budget"`
	Genres              []Genre              `json:"genres,omitempty"`
	Homepage            string               `json:"homepage"`
	Id                  int32                `json:"id"`
	IMDBId              string               `json:"imdb_id"`
	OriginalLanguage    string               `json:"original_language"`
	OriginalTitle       string               `json:"original_title"`
	Overview            string               `json:"overview"`
	Popularity          float64              `json:"popularity"`
	PosterPath          string               `json:"poster_path"`
	ProductionCompanies []ProductionCompany  `json:"production_companies,omitempty"`
	ProductionCountries []ProductionCountry  `json:"production_countries,omitempty"`
	ReleaseDate         string               `json:"release_date"`
	Revenue             int32                `json:"revenue"`
	Runtime             int32                `json:"runtime"`
	SpokenLanguages     []SpokenLanguage     `json:"spoken_languages,omitempty"`
	Status              string               `json:"status"`
	Tagline             string               `json:"tagline"`
	Title               string               `json:"title"`
	Video               bool                 `json:"video"`
	VoteAverage         float64              `json:"vote_average"`
	VoteCount           int32                `json:"vote_count"`
}

type Movies []Movie
