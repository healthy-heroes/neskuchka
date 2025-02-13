import { Provider, defaultTheme, View } from "@adobe/react-spectrum";
import { BrowserRouter, Routes, Route } from "react-router-dom";

import { LandingPage } from "./pages/Landing";

import { NavBar } from "./components/NavBar";

function App() {
  return (
    <Provider theme={defaultTheme}>
      <BrowserRouter>
        <View height="100vh">
          <NavBar />
          <Routes>
            <Route path="/" element={<LandingPage />} />
          </Routes>
        </View>
      </BrowserRouter>
    </Provider>
  );
}

export default App;
