import { LocationProvider, Router, Route } from "preact-iso";
import "./app.css";
import Landing from "./lib/landing";
import Poker from "./lib/poker";
import NotFound from "./lib/notFound";

export function App() {
  return (
    <LocationProvider>
      <Router>
        <Landing path="/" />
        <Poker path="/session/:id" />
        <NotFound default />
      </Router>
    </LocationProvider>
  );
}
