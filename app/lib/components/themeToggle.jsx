import { useContext } from "react";
import { ThemeContext } from "../../app";
import { FiMoon, FiSun } from "react-icons/fi";

export default function ThemeToggle() {
  const { darkMode, setDarkMode } = useContext(ThemeContext);
  function toggleTheme() {
    setDarkMode((prev) => !prev);
  }
  return (
    <li>
      {darkMode ? (
        <a class="secondary" onClick={toggleTheme}>
          <FiMoon size="1.25em" />
        </a>
      ) : (
        <a class="secondary" onClick={toggleTheme}>
          <FiSun size="1.25em" />
        </a>
      )}
    </li>
  );
}
