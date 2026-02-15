// tenant_manager.go

package pkg

// TenantManager handles multi-tenant isolation for the application.
type TenantManager struct {
    tenants map[string]*Tenant
}

// Tenant represents a single tenant in the application.
type Tenant struct {
    ID   string
    Name string
    // other tenant specific fields...
}

// NewTenantManager initializes a new TenantManager.
func NewTenantManager() *TenantManager {
    return &TenantManager{tenants: make(map[string]*Tenant)}
}

// AddTenant adds a new tenant to the management system.
func (tm *TenantManager) AddTenant(id, name string) {
    tm.tenants[id] = &Tenant{ID: id, Name: name}
}

// GetTenant retrieves a tenant by ID.
func (tm *TenantManager) GetTenant(id string) *Tenant {
    return tm.tenants[id]
}

// RemoveTenant removes a tenant from the management system.
func (tm *TenantManager) RemoveTenant(id string) {
    delete(tm.tenants, id)
}