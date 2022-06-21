module go-study/rpc-study/03_service

go 1.18
require geerpc v0.0.2

require codec v0.0.2 // indirect

replace (
	codec v0.0.2 => ./codec
	geerpc v0.0.2 => ./geerpc
)