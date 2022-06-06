module go-study/rpc-study/01_codec

go 1.18

require (
	codec v0.0.1
	geerpc v0.0.1
)

replace (
	codec v0.0.1 => ./codec
	geerpc v0.0.1 => ./geerpc
)
