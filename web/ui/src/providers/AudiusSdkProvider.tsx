import type { AudiusSdk as AudiusSdkType } from "@audius/sdk/dist/sdk/sdk.d.ts";
import type { ServicesConfig } from "@audius/sdk/dist/sdk/config/types.d.ts";
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
  isLoading: boolean;
  isReadOnly: boolean;
};

const AudiusSdkContext = createContext<AudiusSdkContextType>({
  audiusSdk: null,
  isLoading: true,
  isReadOnly: true,
});

export const AudiusSdkProvider = ({ children }: { children: ReactNode }) => {
  const [audiusSdk, setAudiusSdk] = useState<AudiusSdkType | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isReadOnly, setIsReadOnly] = useState(true);
  const envVars = useEnvVars();

  // @ts-expect-error ts(2741). This is only here for debugging and should eventually be removed
  window.audiusSdk = audiusSdk;

  const initSdk = async () => {
    if (!audiusSdk) {
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
      let config = developmentConfig as ServicesConfig;
      let initialSelectedNode = "http://audius-protocol-discovery-provider-1";
      if (envVars.env === "prod") {
        config = productionConfig as ServicesConfig;
        initialSelectedNode = "https://discoveryprovider.audius.co";
      } else if (envVars.env === "stage") {
        config = stagingConfig as ServicesConfig;
        initialSelectedNode = "https://discoveryprovider.staging.audius.co";
      }
      const logger = new Logger({ logLevel: "info" });
      // todo (michelle)
      const apiKey =
        process.env.audius_api_key ||
        "f8f1df516f1ed192c668bf3f781df8db7ed73024";
      const apiSecret =
        process.env.audius_api_secret ||
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
      });
      setAudiusSdk(sdkInst as AudiusSdkType);
    }
    // todo delete read-only?
    setIsReadOnly(false);
    setIsLoading(false);
  };

  useEffect(() => {
    if (window.Web3) {
      void initSdk();
    }
  });

  const contextValue = {
    audiusSdk,
    isLoading,
    isReadOnly,
  };
  return (
    <AudiusSdkContext.Provider value={contextValue}>
      {children}
    </AudiusSdkContext.Provider>
  );
};

// eslint-disable-next-line react-refresh/only-export-components
export const useAudiusSdk = () => useContext(AudiusSdkContext);
