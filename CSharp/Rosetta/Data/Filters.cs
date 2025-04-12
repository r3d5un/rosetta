namespace Rosetta.Data;

public class Filters
{
    public int PageSize { get; set; }
    public Guid? ID { get; set; }
    public string? Name { get; set; }
    public string? Username { get; set; }
    public string? Email { get; set; }
    public DateTime? CreatedAtFrom { get; set; }
    public DateTime? CreatedAtTo { get; set; }
    public DateTime? UpdatedAtFrom { get; set; }
    public DateTime? UpdatedAtTo { get; set; }
    public IList<string> OrderBy { get; set; }
}