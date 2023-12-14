import type { AudiusSdk as AudiusSdkType } from "@audius/sdk/dist/sdk/sdk.d.ts";
import { useAudiusLibs } from "../providers/AudiusLibsProvider";
import {
  ReactNode,
  createContext,
  useContext,
  useState,
  useEffect,
} from "react";
import { useEnvVars } from "../providers/EnvVarsProvider";

type AudiusSdkContextType = {
  audiusSdk: AudiusSdkType | null;
  initSdk: () => Promise<void>;
  removeSdk: () => void;
  isLoading: boolean;
};

const AudiusSdkContext = createContext<AudiusSdkContextType>({
  audiusSdk: null,
  initSdk: async () => {},
  removeSdk: () => {},
  isLoading: true,
});

export const AudiusSdkProvider = ({ children }: { children: ReactNode }) => {
  const { audiusLibs } = useAudiusLibs();
  const [audiusSdk, setAudiusSdk] = useState<AudiusSdkType | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const envVars = useEnvVars();

  // @ts-expect-error ts(2741). This is only here for debugging and should eventually be removed
  window.audiusSdk = audiusSdk;

  const initSdk = async () => {
    console.log("init sdk start");
    if (
      !window.Web3 ||
      !audiusLibs?.Account?.getCurrentUser() ||
      !audiusLibs?.hedgehog
    ) {
      return;
    }

    if (!audiusSdk) {
      // todo (michelle) remove
      console.log("initializing sdk");
      // Dynamically import so sdk uses window.Web3 after it is assigned
      const {
        AppAuth,
        DiscoveryNodeSelector,
        EntityManager,
        Logger,
        StorageNodeSelector,
        developmentConfig,
        stagingConfig,
        productionConfig,
        sdk,
      } = await import("@audius/sdk");
      let config = developmentConfig;
      let initialSelectedNode = "http://audius-protocol-discovery-provider-1";
      if (envVars.env === "prod") {
        config = productionConfig;
        initialSelectedNode = "https://discoveryprovider.audius.co";
      } else if (envVars.env === "stage") {
        config = stagingConfig;
        initialSelectedNode = "https://discoveryprovider.staging.audius.co";
      }
      const logger = new Logger({ logLevel: "info" });
      const apiKey =
        audiusLibs?.hedgehog?.wallet?.getAddressString() ||
        "f8f1df516f1ed192c668bf3f781df8db7ed73024";
      const apiSecret =
        audiusLibs?.hedgehog?.wallet?.getPrivateKeyString() ||
        "b178a83612c99e3ae295743bed7b0186a489cc007985f1a06c6ae873dbdf9145";

      const discoveryNodeSelector = new DiscoveryNodeSelector({
        initialSelectedNode,
      });
      const storageNodeSelector = new StorageNodeSelector({
        auth: new AppAuth(apiKey, apiSecret),
        discoveryNodeSelector: discoveryNodeSelector,
        bootstrapNodes: config.storageNodes,
        logger,
      });
      const sdkInst = sdk({
        services: {
          discoveryNodeSelector,
          entityManager: new EntityManager({
            discoveryNodeSelector,
            web3ProviderUrl: config.web3ProviderUrl,
            contractAddress: config.entityManagerContractAddress,
            identityServiceUrl: config.identityServiceUrl,
            useDiscoveryRelay: true,
            logger,
          }),
          storageNodeSelector,
          logger,
        },
        apiKey: apiKey,
        apiSecret: apiSecret,
        appName: "DDEX Demo",
      });
      setAudiusSdk(sdkInst);
    }
    setIsLoading(false);
  };

  const removeSdk = () => {
    setAudiusSdk(null);
  };

  useEffect(() => {
    void initSdk();
  }, [audiusLibs]);

  const contextValue = {
    audiusSdk,
    initSdk,
    removeSdk,
    isLoading,
  };
  return (
    <AudiusSdkContext.Provider value={contextValue}>
      {children}
    </AudiusSdkContext.Provider>
  );
};

// eslint-disable-next-line react-refresh/only-export-components
export const useAudiusSdk = () => useContext(AudiusSdkContext);
