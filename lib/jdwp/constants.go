package jdwp

import "fmt"

type Error uint16

const None Error = 0                                 //
const InvalidThread Error = 10                       // Passed thread is null, is not a valid thread or has exited.
const InvalidThreadGroup Error = 11                  // Thread group invalid.
const InvalidPriority Error = 12                     // Invalid priority.
const ThreadNotSuspended Error = 13                  // If the specified thread has not been suspended by an event.
const ThreadSuspended Error = 14                     // Thread already suspended.
const ThreadNotAlive Error = 15                      // Thread has not been started or is now dead.
const InvalidObject Error = 20                       // If this reference type has been unloaded and garbage collected.
const InvalidClass Error = 21                        // Invalid class.
const ClassNotPrepared Error = 22                    // Class has been loaded but not yet prepared.
const InvalidMethodid Error = 23                     // Invalid method.
const InvalidLocation Error = 24                     // Invalid location.
const InvalidFieldid Error = 25                      // Invalid field.
const InvalidFrameid Error = 30                      // Invalid jframeID.
const NoMoreFrames Error = 31                        // There are no more Java or JNI frames on the call stack.
const OpaqueFrame Error = 32                         // Information about the frame is not available.
const NotCurrentFrame Error = 33                     // Operation can only be performed on current frame.
const TypeMismatch Error = 34                        // The variable is not an appropriate type for the function used.
const InvalidSlot Error = 35                         // Invalid slot.
const Duplicate Error = 40                           // Item already set.
const NotFound Error = 41                            // Desired element not found.
const InvalidMonitor Error = 50                      // Invalid monitor.
const NotMonitorOwner Error = 51                     // This thread doesn't own the monitor.
const Interrupt Error = 52                           // The call has been interrupted before completion.
const InvalidClassFormat Error = 60                  // The virtual machine attempted to read a class file and determined that the file is malformed or otherwise cannot be interpreted as a class file.
const CircularClassDefinition Error = 61             // A circularity has been detected while initializing a class.
const FailsVerification Error = 62                   // The verifier detected that a class file, though well formed, contained some sort of internal inconsistency or security problem.
const AddMethodNotImplemented Error = 63             // Adding methods has not been implemented.
const SchemaChangeNotImplemented Error = 64          // Schema change has not been implemented.
const InvalidTypestate Error = 65                    // The state of the thread has been modified, and is now inconsistent.
const HierarchyChangeNotImplemented Error = 66       // A direct superclass is different for the new class version, or the set of directly implemented interfaces is different and canUnrestrictedlyRedefineClasses is false.
const DeleteMethodNotImplemented Error = 67          // The new class version does not declare a method declared in the old class version and canUnrestrictedlyRedefineClasses is false.
const UnsupportedVersion Error = 68                  // A class file has a version number not supported by this VM.
const NamesDontMatch Error = 69                      // The class name defined in the new class file is different from the name in the old class object.
const ClassModifiersChangeNotImplemented Error = 70  // The new class version has different modifiers and and canUnrestrictedlyRedefineClasses is false.
const MethodModifiersChangeNotImplemented Error = 71 // A method in the new class version has different modifiers than its counterpart in the old class version and and canUnrestrictedlyRedefineClasses is false.
const NotImplemented Error = 99                      // The functionality is not implemented in this virtual machine.
const NullPointer Error = 100                        // Invalid pointer.
const AbsentInformation Error = 101                  // Desired information is not available.
const InvalidEventType Error = 102                   // The specified event type id is not recognized.
const IllegalArgument Error = 103                    // Illegal argument.
const OutOfMemory Error = 110                        // The function needed to allocate memory and no more memory was available for allocation.
const AccessDenied Error = 111                       // Debugging has not been enabled in this virtual machine. JVMTI cannot be used.
const VmDead Error = 112                             // The virtual machine is not running.
const Internal Error = 113                           // An unexpected internal error has occurred.
const UnattachedThread Error = 115                   // The thread being used to call this function is not attached to the virtual machine. Calls must be made from attached threads.
const InvalidTag Error = 500                         // object type id or class tag.
const AlreadyInvoking Error = 502                    // Previous invoke not complete.
const InvalidIndex Error = 503                       // Index is invalid.
const InvalidLength Error = 504                      // The length is invalid.
const InvalidString Error = 506                      // The string is invalid.
const InvalidClassLoader Error = 507                 // The class loader is invalid.
const InvalidArray Error = 508                       // The array is invalid.
const TransportLoad Error = 509                      // Unable to load the transport.
const TransportInit Error = 510                      // Unable to initialize the transport.
const NativeMethod Error = 511                       //
const InvalidCount Error = 512                       // The count is invalid.

func (e Error) Code() uint16 {
	return uint16(e)
}

func (e Error) String() string {
	return fmt.Sprintf("[%d:%s]", e.Code(), e.Message())
}

func (e Error) Message() string {
	switch e {
	case None:
		return ""
	case InvalidThread:
		return "Passed thread is null, is not a valid thread or has exited."
	case InvalidThreadGroup:
		return "Thread group invalid."
	case InvalidPriority:
		return "Invalid priority."
	case ThreadNotSuspended:
		return "If the specified thread has not been suspended by an event."
	case ThreadSuspended:
		return "Thread already suspended."
	case ThreadNotAlive:
		return "Thread has not been started or is now dead."
	case InvalidObject:
		return "If this reference type has been unloaded and garbage collected."
	case InvalidClass:
		return "Invalid class."
	case ClassNotPrepared:
		return "Class has been loaded but not yet prepared."
	case InvalidMethodid:
		return "Invalid method."
	case InvalidLocation:
		return "Invalid location."
	case InvalidFieldid:
		return "Invalid field."
	case InvalidFrameid:
		return "Invalid jframeID."
	case NoMoreFrames:
		return "There are no more Java or JNI frames on the call stack."
	case OpaqueFrame:
		return "Information about the frame is not available."
	case NotCurrentFrame:
		return "Operation can only be performed on current frame."
	case TypeMismatch:
		return "The variable is not an appropriate type for the function used."
	case InvalidSlot:
		return "Invalid slot."
	case Duplicate:
		return "Item already set."
	case NotFound:
		return "Desired element not found."
	case InvalidMonitor:
		return "Invalid monitor."
	case NotMonitorOwner:
		return "This thread doesn't own the monitor."
	case Interrupt:
		return "The call has been interrupted before completion."
	case InvalidClassFormat:
		return "The virtual machine attempted to read a class file and determined that the file is malformed or otherwise cannot be interpreted as a class file."
	case CircularClassDefinition:
		return "A circularity has been detected while initializing a class."
	case FailsVerification:
		return "The verifier detected that a class file, though well formed, contained some sort of internal inconsistency or security problem."
	case AddMethodNotImplemented:
		return "Adding methods has not been implemented."
	case SchemaChangeNotImplemented:
		return "Schema change has not been implemented."
	case InvalidTypestate:
		return "The state of the thread has been modified, and is now inconsistent."
	case HierarchyChangeNotImplemented:
		return "A direct superclass is different for the new class version, or the set of directly implemented interfaces is different and canUnrestrictedlyRedefineClasses is false."
	case DeleteMethodNotImplemented:
		return "The new class version does not declare a method declared in the old class version and canUnrestrictedlyRedefineClasses is false."
	case UnsupportedVersion:
		return "A class file has a version number not supported by this VM."
	case NamesDontMatch:
		return "The class name defined in the new class file is different from the name in the old class object."
	case ClassModifiersChangeNotImplemented:
		return "The new class version has different modifiers and and canUnrestrictedlyRedefineClasses is false."
	case MethodModifiersChangeNotImplemented:
		return "A method in the new class version has different modifiers than its counterpart in the old class version and and canUnrestrictedlyRedefineClasses is false."
	case NotImplemented:
		return "The functionality is not implemented in this virtual machine."
	case NullPointer:
		return "Invalid pointer."
	case AbsentInformation:
		return "Desired information is not available."
	case InvalidEventType:
		return "The specified event type id is not recognized."
	case IllegalArgument:
		return "Illegal argument."
	case OutOfMemory:
		return "The function needed to allocate memory and no more memory was available for allocation."
	case AccessDenied:
		return "Debugging has not been enabled in this virtual machine. JVMTI cannot be used."
	case VmDead:
		return "The virtual machine is not running."
	case Internal:
		return "An unexpected internal error has occurred."
	case UnattachedThread:
		return "The thread being used to call this function is not attached to the virtual machine. Calls must be made from attached threads."
	case InvalidTag:
		return "object type id or class tag."
	case AlreadyInvoking:
		return "Previous invoke not complete."
	case InvalidIndex:
		return "Index is invalid."
	case InvalidLength:
		return "The length is invalid."
	case InvalidString:
		return "The string is invalid."
	case InvalidClassLoader:
		return "The class loader is invalid."
	case InvalidArray:
		return "The array is invalid."
	case TransportLoad:
		return "Unable to load the transport."
	case TransportInit:
		return "Unable to initialize the transport."
	case NativeMethod:
		return ""
	case InvalidCount:
		return "The count is invalid."
	default:
		return "unknown"
	}
}
