export default function VoteList({ data, votesVisible }) {
  if (!data || !data.length) return <div>Waiting for Voters!</div>;
  const votes = data.map((i) => i.vote).filter((i) => i !== -1);
  const average = votes.reduce((p, c) => p + c, 0) / votes.length;
  const dataSorted = data
    .filter((i) => i.active)
    .sort((a, b) => a.vote - b.vote);
  return (
    <table>
      <thead>
        <tr>
          <th scope="col">User</th>
          <th scope="col">Vote</th>
        </tr>
      </thead>
      <tbody>
        {dataSorted.map(({ name, vote }) => (
          <tr>
            <th scope="row">{name}</th>
            {votesVisible ? (
              <td>{vote === -1 ? "No Vote 🤷" : vote}</td>
            ) : (
              <td>{vote === -1 ? "Waiting for Vote ❓" : "Voted ✅"}</td>
            )}
          </tr>
        ))}
      </tbody>
      <tfoot>
        <tr class="average-row">
          <th scope="row">Average</th>
          <td>{votesVisible && !Number.isNaN(average) ? average : "n/a"}</td>
        </tr>
      </tfoot>
    </table>
  );
}
