import {
  Provider as SpectrumProvider,
  defaultTheme,
  Flex,
} from "@adobe/react-spectrum";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

import queryClient from "./services/api/client";

import { LandingPage } from "./pages/Landing";
import { MainTrackPage } from "./pages/MainTrack";

import { NavBar } from "./components/NavBar";
import { Footer } from "./components/Footer";

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <SpectrumProvider theme={defaultTheme} colorScheme="light" scale="large">
        <BrowserRouter>
          <Flex direction="column" height="100%">
            <NavBar />
            <Routes>
              <Route path="/welcome" element={<LandingPage />} />
              <Route path="/" element={<MainTrackPage />} />
            </Routes>
            <Footer />
          </Flex>
        </BrowserRouter>
      </SpectrumProvider>
      <ReactQueryDevtools />
    </QueryClientProvider>
  );
}

export default App;
