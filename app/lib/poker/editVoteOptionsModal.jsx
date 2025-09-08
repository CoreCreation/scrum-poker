import { useRef, useState } from "preact/hooks";

export default function EditVoteOptionsModal({ open, setOpen, save, current }) {
  const textRef = useRef(null);
  const [valid, setValid] = useState(null);

  function close() {
    setOpen(false);
  }
  function validateList(str) {
    const numbers = str.split(",").map((s) => parseInt(s.trim()));
    return [
      numbers.reduce((p, c) => (p &&= c > 0 && !Number.isNaN(c)), true),
      str,
    ];
  }
  if (!open) {
    return null;
  }
  return (
    <dialog open>
      <article>
        <header>
          <button aria-label="Close" rel="prev" onClick={close}></button>
          <p>
            <strong>Edit Vote Options</strong>
          </p>
        </header>
        <p>A comma separated list of numbers. (e.g. 1, 2, 4, 5, 8)</p>
        <input
          ref={textRef}
          type="text"
          defaultValue={current.join(", ")}
          required="true"
          minLength={1}
          maxLength={64}
          aria-invalid={valid !== null ? !valid : null}
          onInput={() => setValid(textRef.current.checkValidity())}
        ></input>
        <footer>
          <button class="secondary" onClick={close}>
            Cancel
          </button>
          <button
            onClick={() => {
              if (!textRef.current.checkValidity()) return;
              const [valid, numbers] = validateList(textRef.current.value);
              console.log("Dialog is valid:", valid);
              if (!valid) {
                setValid(false);
                return;
              }
              save(numbers);
            }}
          >
            Save
          </button>
        </footer>
      </article>
    </dialog>
  );
}
