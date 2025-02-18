import { Provider, defaultTheme, Flex } from "@adobe/react-spectrum";
import { BrowserRouter, Routes, Route } from "react-router-dom";

import { LandingPage } from "./pages/Landing";
import { MainTrackPage } from "./pages/MainTrack";

import { NavBar } from "./components/NavBar";
import { Footer } from "./components/Footer";

function App() {
  return (
    <Provider theme={defaultTheme} colorScheme="light" scale="large">
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
    </Provider>
  );
}

export default App;
