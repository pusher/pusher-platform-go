// Package pusher provides an interface to interact with the Pusher platform.
// It is mostly intended to be used as a library on top of which products
// that run on the platform are built.
//
// To interact with products that run on the platform, it is better to use the product
// specific library.
//
// This package does not contain any top level exports and is divided into several small
// sub-packages that can be selectively imported to pull in the required functionality.
//
// Interaction with the platform is done by creating Instance objects.
//
// To gain lower level access, a Client maybe constructed and passed into construction
// of an Instance. This allows configuration of the client that is used by the Instance.
// Examples of this can be found in the client package.
package pusher
