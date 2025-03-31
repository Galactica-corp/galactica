protoVer=0.15.3
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

#? proto-gen: Generate Protobuf files
proto-gen:
	@$(protoImage) sh ./scripts/protocgen.sh
	
.PHONY: proto-gen
