import { render, screen } from "@testing-library/react";
import { PitchTable } from "../PitchTable";

describe("PitchTable", () => {
  it("renders 20 data rows plus header", () => {
    render(<PitchTable />);
    expect(screen.getByText("Theme")).toBeInTheDocument();
    const rows = screen.getAllByRole("row");
    expect(rows.length).toBe(21);
  });
});
