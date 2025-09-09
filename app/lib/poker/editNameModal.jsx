import { useRef, useState } from "preact/hooks";

export default function EditNameModal({ open, setOpen, save, name }) {
  const textRef = useRef(null);
  const [valid, setValid] = useState(null);

  function close() {
    setOpen(false);
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
            <strong>Edit Username</strong>
          </p>
        </header>
        <input
          ref={textRef}
          type="text"
          defaultValue={name}
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
            onClick={() =>
              textRef.current.checkValidity() && save(textRef.current.value)
            }
          >
            Save
          </button>
        </footer>
      </article>
    </dialog>
  );
}
