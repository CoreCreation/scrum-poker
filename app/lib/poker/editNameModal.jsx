import { useRef } from "preact/hooks";

export default function EditNameModal({ open, save }) {
  const textRef = useRef(null);
  return (
    <dialog open={open}>
      <input ref={textRef} type="text"></input>
      <button onClick={() => save(textRef.current.value)}>Save</button>
    </dialog>
  );
}
