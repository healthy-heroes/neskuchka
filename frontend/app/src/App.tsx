import { Provider, defaultTheme, Flex, Content } from "@adobe/react-spectrum";
import { BrowserRouter, Routes, Route } from "react-router-dom";

import { LandingPage } from "./pages/Landing";

import { NavBar } from "./components/NavBar";
import { Footer } from "./components/Footer";

function App() {
  return (
    <Provider theme={defaultTheme} colorScheme="light" scale="large">
      <BrowserRouter>
        <Flex direction="column" height="100%">
          <NavBar />
          <Content>
            <Routes>
              <Route path="/welcome" element={<LandingPage />} />
            </Routes>
          </Content>
          <Footer />
        </Flex>
      </BrowserRouter>
    </Provider>
  );
}

export default App;
