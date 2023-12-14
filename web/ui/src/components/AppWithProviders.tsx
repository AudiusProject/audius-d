import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { WagmiProvider } from "wagmi";
import { AudiusLibsProvider } from "../providers/AudiusLibsProvider.tsx";
import { AudiusSdkProvider } from "../providers/AudiusSdkProvider"
import App from "./App.tsx";
import { useWagmiConfig } from "../hooks/useWagmiConfig.tsx";
import Web3 from "web3";
import type Web3Type from "web3";

declare global {
  interface Window {
    Web3: Web3Type
  }
}

const AppWithProviders = () => {
  window.Web3 = Web3;
  const wagmiConfig = useWagmiConfig();

  const queryClient = new QueryClient();

  return (
    <WagmiProvider config={wagmiConfig}>
      <QueryClientProvider client={queryClient}>
        <AudiusLibsProvider>
          <AudiusSdkProvider>
            <App />
          </AudiusSdkProvider>
        </AudiusLibsProvider>
        <ReactQueryDevtools initialIsOpen={false} />
      </QueryClientProvider>
    </WagmiProvider>
  );
};

export default AppWithProviders;
