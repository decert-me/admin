import './styles/App.css';
import './styles/index.scss';
import BeforeRouterEnter from './components/BeforeRouterEnter';
import { WagmiConfig, configureChains, createConfig } from 'wagmi';
import { polygon, polygonMumbai } from 'viem/chains';
import { MetaMaskConnector } from 'wagmi/connectors/metaMask';
import { WalletConnectLegacyConnector } from 'wagmi/connectors/walletConnectLegacy';
// import { PhantomConnector } from 'phantom-wagmi-connector';

import { infuraProvider } from 'wagmi/providers/infura';
import { publicProvider } from 'wagmi/providers/public';


const { chains, publicClient, webSocketPublicClient } = configureChains(
  [polygon, polygonMumbai],
  [
    infuraProvider({ apiKey: process.env.REACT_APP_INFURA_API_KEY }),
    publicProvider(),
  ],
)

const config = createConfig({
  autoConnect: true,
  connectors: [
    new MetaMaskConnector({ chains }),
    // new PhantomConnector({ chains }),
    // new WalletConnectConnector({
    //   chains,
    //   options: {
    //     projectId: process.env.NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID ?? '',
    //   },
    // }),
  ],
  publicClient,
  webSocketPublicClient,
})


function App() {
  window.Buffer = window.Buffer || require("buffer").Buffer;

  return (
    <WagmiConfig config={config}>
        <BeforeRouterEnter />
    </WagmiConfig>
  );
}

export default App;
