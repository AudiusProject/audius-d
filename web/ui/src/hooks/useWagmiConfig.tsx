import type { Transport } from "viem";
import { useMemo } from "react";
import { http, createConfig, fallback, webSocket } from "wagmi";
import { mainnet, goerli } from "wagmi/chains";
import { useEnvVars } from "../providers/EnvVarsProvider.tsx";
import { ethers } from "ethers";
// import { metaMask } from 'wagmi/connectors'

export const useWagmiConfig = () => {
  const { ethProviderUrl } = useEnvVars();
  const localEndpoints = [
    "http://audius-protocol-eth-ganache-1",
    "http://localhost:8546",
  ];

  return useMemo(() => {
    // audius-docker-compose configs only allow for one RPC env var, which could be a single endpoint or a comma-separated list of endpoint
    const providerEndpoints = ethProviderUrl.includes(",")
      ? ethProviderUrl.split(",")
      : [ethProviderUrl];

    providerEndpoints.forEach((url: string) => {
      if (localEndpoints.includes(url)) {
        const provider = new ethers.providers.JsonRpcProvider(url);
        // advance local ganache chain
        provider.send("evm_mine", []);
      }
    });

    const rpcProviders: Transport[] = providerEndpoints.map((url: string) =>
      url.startsWith("ws") ? webSocket(url) : http(url),
    );

    // Allows for fallback to other RPC endpoints if the first one fails
    const transports = fallback([
      ...rpcProviders,
      http(localEndpoints[0]),
      http(), // Public fallback provider (rate-limited)
    ]);

    const wagmiConfig = createConfig({
      chains: [mainnet, goerli],
      // connectors: [metaMask()],
      transports: {
        [mainnet.id]: transports,
        [goerli.id]: transports,
      },
    });

    return wagmiConfig;
  }, [ethProviderUrl]);
};

declare module "wagmi" {
  interface Register {
    config: ReturnType<typeof useWagmiConfig>;
  }
}
