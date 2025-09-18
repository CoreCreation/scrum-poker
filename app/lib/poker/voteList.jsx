export default function VoteList({ data, votesVisible }) {
  if (!data || !data.length) return <div>Waiting for Voters!</div>;
  const votes = data
    .filter((i) => i.active)
    .map((i) => i.vote)
    .filter((i) => i !== -1);
  const average = votes.reduce((p, c) => p + c, 0) / votes.length;
  const dataSorted = data.sort((a, b) => a.uuid.localeCompare(b.uuid));
  return (
    <table>
      <thead>
        <tr>
          <th scope="col">User</th>
          <th scope="col">Vote</th>
        </tr>
      </thead>
      <tbody>
        {dataSorted.map(({ uuid, username, vote }) => (
          <tr id={uuid}>
            <th scope="row">{username.length ? username : "No Name"}</th>
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
