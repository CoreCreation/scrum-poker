import { LocationProvider, Router } from "preact-iso";
import "./app.css";
import Landing from "./lib/landing";
import Poker from "./lib/poker";
import NotFound from "./lib/notFound";
import { createContext, useEffect } from "react";
import { useRef, useState } from "preact/hooks";

export const ThemeContext = createContext();

export function App() {
  const [darkMode, setDarkMode] = useState(
    localStorage.getItem("scrum-poker-dark-mode") ??
      window.matchMedia("(prefers-color-scheme: dark)").matches
  );

  useEffect(() => {
    htmlRef.current.dataset.theme = darkMode ? "dark" : "light";
  }, [darkMode]);

  const htmlRef = useRef(document.querySelector("html"));
  return (
    <ThemeContext.Provider value={{ darkMode, setDarkMode }}>
      <LocationProvider>
        <Router>
          <Landing path="/" />
          <Poker path="/session/:id" />
          <NotFound default />
        </Router>
      </LocationProvider>
    </ThemeContext.Provider>
  );
}
