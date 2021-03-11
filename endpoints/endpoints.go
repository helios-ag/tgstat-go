package endpoints

// API endpoints
const (
	ChannelsGet         string = "/channels/get"
	ChannelsSearch      string = "/channels/search"
	ChannelsStat        string = "/channels/stat"
	ChannelsMention     string  = "/channels/mentions"
	ChannelsForward     string = "/channels/forwards"
	ChannelsSubscribers string = "/channels/subscribers"
	ChannelsView        string = "/channels/views"
	ChannelsAdd         string = "/channels/add"

	PostsGet string = "/posts/get"
	PostsStat string = "/posts/stat"
	PostsSearch string = "/posts/search"

	WordsMentionsByPeriod string = "/words/mentions-by-period"
	WordsMentionsByChannels string = "/words/mentions-by-channels"

	UsageStat string = "/usage/stat"

	DatabaseCategories string = "/database/categories"
	DatabaseCountries string = "/database/countries"
	DatabaseLanguages string = "/database/languages"
)

