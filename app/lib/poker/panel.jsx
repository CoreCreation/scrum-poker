import VoteButton from "./panel/voteButton";

export default function Panel({ options }) {
  return (
    <div>
      {options.map((option) => (
        <VoteButton option={option} />
      ))}
    </div>
  );
}
