import { AudiusLibs, sdk, AudiusSdk, stagingConfig } from "@audius/sdk";

// returns authed sdk instance
export const initSdkUser = async (): Promise<AudiusSdk> => {
  const web3Config = {};
  const ethWeb3Config = {};
  const solanaWeb3Config = {};
  const identityServiceConfig = {
    url: stagingConfig.identityServiceUrl,
    useHedgehogLocalStorage: false,
  };
  const discoveryProviderConfig = {};
  const creatorNodeConfig = {};
  const libs = new AudiusLibs({
    web3Config,
    ethWeb3Config,
    solanaWeb3Config,
    discoveryProviderConfig,
    identityServiceConfig,
    creatorNodeConfig,
    isServer: true,
    isDebug: true,
    preferHigherPatchForPrimary: true,
    preferHigherPatchForSecondaries: true,
    useDiscoveryRelay: true,
    logger: console,
  });

  return sdk({ apiKey: "", apiSecret: "" });
};
