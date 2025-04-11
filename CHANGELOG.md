# Axone SDK changelog

## [1.2.0](https://github.com/axone-protocol/axone-sdk/compare/v1.1.0...v1.2.0) (2025-04-11)


### Features

* **dataverse:** add support for retrieving cognitarium info ([1b75d45](https://github.com/axone-protocol/axone-sdk/commit/1b75d45a7d59566f67bd25436ca66579c2cdf8bd))
* **dataverse:** add support for retrieving dataverse info ([a0b4896](https://github.com/axone-protocol/axone-sdk/commit/a0b48961ba3171491786817635bbdace703bb42c))

## [1.1.0](https://github.com/axone-protocol/axone-sdk/compare/v1.0.0...v1.1.0) (2025-04-05)


### Features

* **dataverse:** add GovCode to retrieve governance code from law-stone ([09ae724](https://github.com/axone-protocol/axone-sdk/commit/09ae724daf3fb3f854273c308bf757162434baec))


### Bug Fixes

* **srv:** avoid unnecessary HTTP 200 status write ([02256d7](https://github.com/axone-protocol/axone-sdk/commit/02256d7b2502a54497ebb94469f5f9e4c8f8c69d))

## 1.0.0 (2024-11-13)


### Features

* add simple cryptographic key management ([42b57d4](https://github.com/axone-protocol/axone-sdk/commit/42b57d4bbb28b8018c2054146b9ea027d4f36b25))
* add storage proxy http configuration ([f5736f6](https://github.com/axone-protocol/axone-sdk/commit/f5736f6002f58e876c56b66fa90c001026a67655))
* **auth:** add check on credential parsing ([899ff8a](https://github.com/axone-protocol/axone-sdk/commit/899ff8a76ae38ab84582183db7334c14d473b779))
* **auth:** implemant auth handler ([8574af4](https://github.com/axone-protocol/axone-sdk/commit/8574af42f9384dac169c65b33d9e38612255720b))
* **auth:** parse auth claim credentials ([ddc64e7](https://github.com/axone-protocol/axone-sdk/commit/ddc64e71330f2e833d52d44bbce84cac2893a446))
* **credential:** add a secp256k1 public key fetcher for parse credentials ([b03b381](https://github.com/axone-protocol/axone-sdk/commit/b03b381d7b3d98ebdfb788d709424ccbbc0eed45))
* **credential:** add better error handling ([748e7a4](https://github.com/axone-protocol/axone-sdk/commit/748e7a44e55aea13ec92e9b9c44cbbdb382769dd))
* **credential:** add method to parse vc ([4427d70](https://github.com/axone-protocol/axone-sdk/commit/4427d70443b39b2fc6954d5deabdfae7d57fbadf))
* **credential:** better manage error typed ([d46baf3](https://github.com/axone-protocol/axone-sdk/commit/d46baf30837c6253497e058390ccadd284cf75e0))
* **credential:** check authentication proof purpose for auth claim ([28a2844](https://github.com/axone-protocol/axone-sdk/commit/28a2844fea9b61a060b4122f2d7a953164759696))
* **credential:** create a generator allowing generated credential ([2b29672](https://github.com/axone-protocol/axone-sdk/commit/2b29672560752fd8e6ba19fa569061b91d9e4a28))
* **credential:** create dataset description credential ([c7d32af](https://github.com/axone-protocol/axone-sdk/commit/c7d32aff9290e7740e5b381c0de7582934a25a45))
* **credential:** create publication credential ([15b3358](https://github.com/axone-protocol/axone-sdk/commit/15b3358c736097be593cab3c996d77f88cf8417b))
* **credential:** handle future issued credentials ([be4aecd](https://github.com/axone-protocol/axone-sdk/commit/be4aecd892394203ce81c456c51797476f11fbbd))
* **credential:** put Parse method public ([1593420](https://github.com/axone-protocol/axone-sdk/commit/15934206b1e03808fb20c7996dcce89e7420a9a5))
* **credential:** use keyring for signing credential ([5b391c0](https://github.com/axone-protocol/axone-sdk/commit/5b391c03f78605f5d7957200b8cae94c7c5ff4ff))
* **dataverse:** add GetGovCode interface ([6318ae8](https://github.com/axone-protocol/axone-sdk/commit/6318ae8790375270183d1a5d40dd9481e7237f44))
* **dataverse:** convert credential to rdf ([e4d308d](https://github.com/axone-protocol/axone-sdk/commit/e4d308d1a31207f136ffc0934f4503aa974a6ee7))
* **dataverse:** get cognitarium address at NewClient ([7c33523](https://github.com/axone-protocol/axone-sdk/commit/7c33523cfe248bc7eb5b894ab4ae32e04c0b5957))
* **dataverse:** handle more properly errors ([e89c911](https://github.com/axone-protocol/axone-sdk/commit/e89c911014b89c957fe8cd8a6cd9c367c0b0babc))
* **dataverse:** implement client with GetGov addresse ([678df83](https://github.com/axone-protocol/axone-sdk/commit/678df83cd423409c20d7269549f3d4849dbf892a))
* **dataverse:** implement get resource gov addr ([2bd6f53](https://github.com/axone-protocol/axone-sdk/commit/2bd6f5346adfb52e4f067a0b9faabe5763b6f1af))
* **dataverse:** submit claims ([473d562](https://github.com/axone-protocol/axone-sdk/commit/473d562751c4d783a53061188efbfa672c04788f))
* **dataverse:** submit claims return tx response ([02e1a67](https://github.com/axone-protocol/axone-sdk/commit/02e1a67f3f77a9d9b56ff9ce8ae0379ac4c086a9))
* design a http server with easy route configuration ([7b7dd89](https://github.com/axone-protocol/axone-sdk/commit/7b7dd892e9a1a5f856c3b92eb0eca979dd5af3a0))
* design authentication proxy elements ([90c79be](https://github.com/axone-protocol/axone-sdk/commit/90c79beac3962c5f5ae866702a4c6a41e6e13f4e))
* design storage service proxy ([cd4d012](https://github.com/axone-protocol/axone-sdk/commit/cd4d0124f9975403781c9c8e6567ed2beddedb94))
* **http:** add listen method to http server ([74fceda](https://github.com/axone-protocol/axone-sdk/commit/74fceda84c42c6b33f385c1a14f07a94e4f36532))
* **http:** ease server configuration ([aa419e8](https://github.com/axone-protocol/axone-sdk/commit/aa419e8d7951e858ea66427b5a3d9b39e712297b))
* implement goverance checks against law stone ([7908cac](https://github.com/axone-protocol/axone-sdk/commit/7908cacc3b019c1597928ef272c32e05857fb452))
* **keys:** incldue addr, didKeyid on Keyring interface ([d972a87](https://github.com/axone-protocol/axone-sdk/commit/d972a879482547148bd8800b2793d7c7dca3f85b))
* pre implement jwt management ([b531ed6](https://github.com/axone-protocol/axone-sdk/commit/b531ed6155af46ea7089cd4ccce6782ad4385e00))
* properly manage http responses ([65d8de2](https://github.com/axone-protocol/axone-sdk/commit/65d8de20b5a5de9b0c93d717eab5af5c36ab33e9))
* setup dataverse client dummy interface ([51347fa](https://github.com/axone-protocol/axone-sdk/commit/51347fa5608054b133e442f593f32a7e575a1c66))
* **tx:** add tx client to send transaction ([59becfb](https://github.com/axone-protocol/axone-sdk/commit/59becfb84dcc6ecdd852f24f6bb1af30cd7ea24a))


### Bug Fixes

* **auth:** change status code returned on invalid body read ([4102bc2](https://github.com/axone-protocol/axone-sdk/commit/4102bc2f93345fab390e4f955db771cf589284aa))
* **auth:** check read access with provided identity ([6342a38](https://github.com/axone-protocol/axone-sdk/commit/6342a384626dc66829d778d01d4396972c9fb168))
* **credential:** change governance proof purpose as assertionMethod ([672fc85](https://github.com/axone-protocol/axone-sdk/commit/672fc85aaed24a40a5aa3b5bad5de818aca0a902))
* **dataverse:** allow getting governance addr in URI ([f041b75](https://github.com/axone-protocol/axone-sdk/commit/f041b7504c5761c3fc36f64297b488ab69d2dbfe))
* **dataverse:** remove the unused cognitarium address hold by client ([bc5b3ad](https://github.com/axone-protocol/axone-sdk/commit/bc5b3ad9f6dc16b2376d389c6c828ab7b75429c5))
* **dataverse:** use given context for fetch cognitarium ([92cbf72](https://github.com/axone-protocol/axone-sdk/commit/92cbf72486a6dbc3a880e99bbdc8f7a687daaeba))
* **http:** expose authenticate http endpoint on POST req ([e4e82b8](https://github.com/axone-protocol/axone-sdk/commit/e4e82b8b432eceb4659b0d132ee3f045566cfe6a))
